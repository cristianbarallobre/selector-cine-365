package main

// Shared TMDB identity catalogue for all collections. Availability and ratings
// remain separate: two titles can represent the same film even when neither
// is currently available on a Spanish streaming service.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type filmIdentityRecord struct {
	TMDBID  int    `json:"tmdb_id"`
	Updated string `json:"updated"`
	State   string `json:"state"` // ok or not-found
}
type identityStorage struct {
	Entries map[string]filmIdentityRecord `json:"entries"`
}

// Serialize requests for new IDs to avoid simultaneous duplicate lookups.
var identityRequestLock sync.Mutex

func (p *availabilityServer) saveIdentityCacheLocked() {
	if p.identityPath == "" {
		return
	}
	bytes, err := json.MarshalIndent(identityStorage{Entries: p.identities}, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(p.identityPath), 0700)
	tmp := p.identityPath + ".tmp"
	if os.WriteFile(tmp, bytes, 0600) == nil {
		if os.Rename(tmp, p.identityPath) != nil {
			_ = os.WriteFile(p.identityPath, bytes, 0600)
			_ = os.Remove(tmp)
		}
	}
}

func (p *availabilityServer) knownIdentities() map[string]int {
	p.mu.Lock()
	defer p.mu.Unlock()
	ids := map[string]int{}
	for key, record := range p.identities {
		if record.TMDBID > 0 {
			ids[key] = record.TMDBID
		}
	}
	for key, record := range p.cache {
		if record.TMDBID > 0 {
			ids[key] = record.TMDBID
		}
	}
	// A manual correction is always authoritative, including over old cache.
	for key, id := range p.overrides {
		if id > 0 {
			ids[key] = id
		}
	}
	return ids
}
func (p *availabilityServer) knownIdentity(key string) (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if id := p.overrides[key]; id > 0 {
		return id, true
	}
	if entry := p.identities[key]; entry.TMDBID > 0 {
		return entry.TMDBID, true
	}
	if item := p.cache[key]; item.TMDBID > 0 {
		return item.TMDBID, true
	}
	if entry, found := p.identities[key]; found && entry.State == "not-found" {
		stamp, e := time.Parse(time.RFC3339, entry.Updated)
		if e == nil && time.Since(stamp) < 48*time.Hour {
			return 0, true
		}
	}
	return 0, false
}
func (p *availabilityServer) recordIdentity(key string, id int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.identities == nil {
		p.identities = map[string]filmIdentityRecord{}
	}
	if id > 0 {
		p.identities[key] = filmIdentityRecord{TMDBID: id, State: "ok", Updated: time.Now().Format(time.RFC3339)}
	} else {
		p.identities[key] = filmIdentityRecord{TMDBID: 0, State: "not-found", Updated: time.Now().Format(time.RFC3339)}
	}
	p.saveIdentityCacheLocked()
}

// Identify only: it does not query availability, film credits or posters.
// This allows previously saved lists to be reconciled without making thousands
// of additional /watch/providers requests.
func (p *availabilityServer) identifyFilm(ctx context.Context, title string, year int, force bool) (int, string, error) {
	key := movieCacheKey(title, year)
	if !force {
		if id, ok := p.knownIdentity(key); ok {
			if id > 0 {
				return id, "ok", nil
			}
			return 0, "not-found", nil
		}
	}
	if p.getKey() == "" {
		return 0, "unconfigured", nil
	}
	identityRequestLock.Lock()
	defer identityRequestLock.Unlock()
	if !force {
		if id, ok := p.knownIdentity(key); ok {
			if id > 0 {
				return id, "ok", nil
			}
			return 0, "not-found", nil
		}
	}
	type film struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Original string `json:"original_title"`
		Released string `json:"release_date"`
	}
	var id int
	for _, query := range titleSearchVariants(title) {
		var search struct {
			Results []film `json:"results"`
		}
		params := url.Values{"language": {"es-ES"}, "query": {query}, "include_adult": {"false"}}
		if year > 0 {
			params.Set("year", strconv.Itoa(year))
		}
		if err := requestTMDB(ctx, "/search/movie", params, p.getKey(), &search); err != nil {
			return 0, "error", err
		}
		// If the year-specific search has candidates but none match securely,
		// retry WITHOUT year. Some foreign releases differ by a year or two.
		candidates := func(films []film) []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Original string `json:"original_title"`
			Released string `json:"release_date"`
		} {
			out := make([]struct {
				ID       int    `json:"id"`
				Title    string `json:"title"`
				Original string `json:"original_title"`
				Released string `json:"release_date"`
			}, 0, len(films))
			for _, m := range films {
				out = append(out, struct {
					ID       int    `json:"id"`
					Title    string `json:"title"`
					Original string `json:"original_title"`
					Released string `json:"release_date"`
				}{m.ID, m.Title, m.Original, m.Released})
			}
			return out
		}
		candidate, _ := chooseMovie(query, year, candidates(search.Results))
		if candidate == 0 && year > 0 {
			params.Del("year")
			if err := requestTMDB(ctx, "/search/movie", params, p.getKey(), &search); err != nil {
				return 0, "error", err
			}
			candidate, _ = chooseMovie(query, year, candidates(search.Results))
		}
		if candidate > 0 {
			// A film without a year must be highly unambiguous: exact title or
			// original-title match. Otherwise prefer an unresolved film to a
			// false positive that might sync the wrong watched status.
			if year == 0 {
				exactMatches := map[int]bool{}
				for _, m := range search.Results {
					if normaliseTitle(m.Title) == normaliseTitle(query) || normaliseTitle(m.Original) == normaliseTitle(query) {
						exactMatches[m.ID] = true
					}
				}
				// Several different films have this same name: request a
				// year or a manual ID instead of merging unrelated titles.
				if len(exactMatches) != 1 || !exactMatches[candidate] {
					continue
				}
			}
			id = candidate
			break
		}
	}
	p.recordIdentity(key, id)
	if id == 0 {
		return 0, "not-found", nil
	}
	return id, "ok", nil
}
func (a *runningApp) serveIdentitiesHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	switch r.URL.Path {
	case "/api/identities":
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"configured": a.platforms.configured(), "entries": a.platforms.knownIdentities()})
	case "/api/identities/resolve":
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var request struct {
			Title string `json:"title"`
			Year  int    `json:"year"`
			Force bool   `json:"force"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&request) != nil {
			http.Error(w, "Datos incorrectos", 400)
			return
		}
		request.Title = strings.TrimSpace(request.Title)
		if len(request.Title) == 0 || len(request.Title) > 250 || request.Year < 0 || request.Year > 2200 {
			http.Error(w, "Película inválida", 400)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 40*time.Second)
		defer cancel()
		id, state, err := a.platforms.identifyFilm(ctx, request.Title, request.Year, request.Force)
		if err != nil {
			var note = "Error consultando TMDB: " + err.Error()
			// Return a well-formed response: the UI can stop or retry.
			_ = json.NewEncoder(w).Encode(map[string]any{"tmdb_id": 0, "state": "error", "note": note})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"tmdb_id": id, "state": state})
	default:
		http.Error(w, fmt.Sprintf("Ruta desconocida: %s", r.URL.Path), 404)
	}
}
