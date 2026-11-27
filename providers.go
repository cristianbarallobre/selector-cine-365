package main

// Consulta de plataformas en España. Fuente: datos de JustWatch distribuidos
// por TMDB. Solo el usuario proporciona su propia clave de API de TMDB.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

var tmdbClient = &http.Client{Timeout: 15 * time.Second}

type platformSettings struct {
	Key string `json:"api_key,omitempty"`
 Protected []byte `json:"protected,omitempty"`
 Keychain bool `json:"keychain,omitempty"`
}
type providerOffer struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}
type providerResult struct {
	State           string   `json:"state"`
	Title           string   `json:"title,omitempty"`
	FilmURL         string   `json:"film_url,omitempty"`
	Subscriptions   []string `json:"subscriptions"`
	Free            []string `json:"free"`
	Ads             []string `json:"ads"`
	Rent            []string `json:"rent"`
	Buy             []string `json:"buy"`
	Date            string   `json:"date,omitempty"`
	Note            string   `json:"note,omitempty"`
	VoteAverage     float64  `json:"vote_average,omitempty"`
	VoteCount       int      `json:"vote_count,omitempty"`
	PosterPath      string   `json:"poster_path,omitempty"`
	Genres          []string `json:"genres,omitempty"`
	Countries       []string `json:"countries,omitempty"`
	Director        string   `json:"director,omitempty"`
	Runtime         int      `json:"runtime,omitempty"`
	MetadataVersion int      `json:"metadata_version,omitempty"`
	TMDBID          int      `json:"tmdb_id,omitempty"`
}
type platformCache struct {
	Entries map[string]providerResult `json:"entries"`
}
type availabilityServer struct {
	mu               sync.Mutex
	settingsPath     string
	cachePath        string
	posterDir        string
	legacyPosterDir  string
	overridesPath    string
	overrides        map[string]int
	identityPath     string
	identities       map[string]filmIdentityRecord
	apiKey           string
	cache            map[string]providerResult
	cacheDirty       bool
	identityDirty    bool
	cacheTimer       *time.Timer
	overrideRevision map[string]uint64
}

func newAvailabilityServer(base string, cacheDirs ...string) *availabilityServer {
	cacheDir := base
	if len(cacheDirs) > 0 && cacheDirs[0] != "" {
		cacheDir = cacheDirs[0]
	}
	p := &availabilityServer{settingsPath: filepath.Join(base, "plataformas_ajustes.json"), cachePath: filepath.Join(cacheDir, "plataformas_cache.json"), posterDir: filepath.Join(cacheDir, "carteles"), overridesPath: filepath.Join(base, "identificaciones_tmdb.json"), overrides: make(map[string]int), identityPath: filepath.Join(base, "identidades_cine.json"), identities: make(map[string]filmIdentityRecord), cache: make(map[string]providerResult)}
	if cacheDir != base {
		p.legacyPosterDir = filepath.Join(base, "carteles")
		_ = copyLegacyCacheIfMissing(filepath.Join(base, "plataformas_cache.json"), p.cachePath, 12<<20)
		_, _ = prunePosterCache(p.posterDir, appPosterMaxFiles, appPosterMaxBytes, appPosterMaxAge)
	}
	if key,e:=readTMDBKey(p.settingsPath);e==nil{p.apiKey=key}else if !os.IsNotExist(e){diagnosticError("tmdb","leer credencial local",e)}
	b, e := os.ReadFile(p.cachePath)
	if e == nil {
		var c platformCache
		if json.Unmarshal(b, &c) == nil && c.Entries != nil {
			p.cache = c.Entries
			if trimmed, err := p.encodeProviderCacheLocked(); err == nil {
				// Mantener la copia local ajustada incluso antes de recibir nuevas fichas.
				if len(trimmed) < len(b) {
					if err := os.MkdirAll(filepath.Dir(p.cachePath), 0700); err == nil {
						_ = replaceLibraryFile(p.cachePath, trimmed)
					}
				}
			}
		}
	}
	if data, e := os.ReadFile(p.overridesPath); e == nil {
		var overrides map[string]int
		if json.Unmarshal(data, &overrides) == nil && overrides != nil {
			p.overrides = overrides
		}
	}
	if data, e := os.ReadFile(p.identityPath); e == nil {
		var saved identityStorage
		if json.Unmarshal(data, &saved) == nil && saved.Entries != nil {
			p.identities = saved.Entries
		}
	}
	return p
}
func (p *availabilityServer) configured() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.apiKey != ""
}
func (p *availabilityServer) getKey() string { p.mu.Lock(); defer p.mu.Unlock(); return p.apiKey }
func (p *availabilityServer) saveKey(key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(p.settingsPath), 0700); err != nil {
		return err
	}
	if e:=writeTMDBKey(p.settingsPath,key);e!=nil{return e}
	p.apiKey = key
	return nil
}

// No se comparten las credenciales API; solo IDs confirmados por el usuario.
func (p *availabilityServer) manualOverrides() map[string]int {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make(map[string]int, len(p.overrides))
	for k, id := range p.overrides {
		result[k] = id
	}
	return result
}
func (p *availabilityServer) importManualOverrides(mappings map[string]int) error {
	if len(mappings) > 10000 {
		return errors.New("Demasiadas identificaciones")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	next := make(map[string]int, len(p.overrides)+len(mappings))
	for k, id := range p.overrides {
		next[k] = id
	}
	changed := make(map[string]int)
	for key, id := range mappings {
		if len(key) < 3 || len(key) > 275 || !strings.Contains(key, "|") || id < 0 || id > 999999999 {
			return errors.New("Identificación TMDB incorrecta")
		}
		old, exists := next[key]
		if id == 0 {
			_, hasIdentity := p.identities[key]
			_, hasAvailability := p.cache[key]
			if exists || hasIdentity || hasAvailability {
				delete(next, key)
				changed[key] = 0
			}
		} else if !exists || old != id {
			next[key] = id
			changed[key] = id
		}
	}
	if len(changed) == 0 {
		return nil
	}
	bytes, e := json.MarshalIndent(next, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(p.overridesPath), 0700); e != nil {
		return e
	}
	if e = replaceLibraryFile(p.overridesPath, bytes); e != nil {
		return e
	}
	p.overrides = next
	stamp := time.Now().Format(time.RFC3339)
	for key, id := range changed {
		p.bumpOverrideRevisionLocked(key)
		if id == 0 {
			delete(p.identities, key)
		} else {
			p.identities[key] = filmIdentityRecord{TMDBID: id, State: "ok", Updated: stamp}
		}
		if id == 0 || p.cache[key].TMDBID != id {
			delete(p.cache, key)
		}
	}
	p.identityDirty = true
	p.cacheDirty = true
	p.flushProviderCacheLocked()
	return nil
}
func (p *availabilityServer) getCached(k string, force bool) (providerResult, bool) {
	if force {
		return providerResult{}, false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	r, ok := p.cache[k]
	if !ok {
		return r, false
	}
	if manual := p.overrides[k]; manual > 0 && manual != r.TMDBID {
		return providerResult{}, false
	}
	// Los datos se conservan para uso offline; la consulta comprueba su caducidad.
	if r.MetadataVersion < 3 {
		return providerResult{}, false
	}
	return r, true
}
func (p *availabilityServer) bumpOverrideRevisionLocked(key string) {
	if p.overrideRevision == nil {
		p.overrideRevision = make(map[string]uint64)
	}
	p.overrideRevision[key]++
}
func (p *availabilityServer) cachedEntries() map[string]providerResult {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make(map[string]providerResult, len(p.cache))
	for k, v := range p.cache {
		if manual := p.overrides[k]; manual > 0 && manual != v.TMDBID {
			continue
		}
		result[k] = v
	}
	return result
}

// Mantener un catálogo acotado. Las identificaciones confirmadas viven en
// identidades_cine.json, separado de esta caché regenerable.
func (p *availabilityServer) encodeProviderCacheLocked() ([]byte, error) {
	const maxEntries = 3200
	const maxJSONBytes = 10 << 20
	data, err := json.Marshal(platformCache{Entries: p.cache})
	if err != nil || (len(p.cache) <= maxEntries && len(data) <= maxJSONBytes) {
		return data, err
	}
	type entry struct {
		key  string
		date time.Time
	}
	rows := make([]entry, 0, len(p.cache))
	for key, record := range p.cache {
		date, _ := time.Parse(time.RFC3339, record.Date)
		rows = append(rows, entry{key, date})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].date.Equal(rows[j].date) {
			return rows[i].key < rows[j].key
		}
		return rows[i].date.Before(rows[j].date)
	})
	size := len(data)
	for _, row := range rows {
		if len(p.cache) <= maxEntries && size <= maxJSONBytes {
			break
		}
		key, _ := json.Marshal(row.key)
		record, _ := json.Marshal(p.cache[row.key])
		size -= len(key) + len(record) + 1 // Un separador adicional puede sobrar.
		delete(p.cache, row.key)
	}
	return json.Marshal(platformCache{Entries: p.cache})
}
func (p *availabilityServer) putCached(k string, r providerResult) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.putCachedLocked(k, r)
}

func (p *availabilityServer) finishAvailability(key string, revision uint64, result providerResult) providerResult {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.overrideRevision[key] != revision {
		if current, ok := p.cache[key]; ok && current.MetadataVersion >= 3 && (p.overrides[key] == 0 || p.overrides[key] == current.TMDBID) {
			return current
		}
		return providerResult{State: "superseded"}
	}
	if prior := p.cache[key]; prior.TMDBID > 0 && (prior.State=="ok"||prior.State=="no-data") &&
       (p.overrides[key]==0||p.overrides[key]==prior.TMDBID) && (result.TMDBID==0||result.TMDBID==prior.TMDBID) && (result.State=="error"||result.State=="not-found") {return prior}
	if prior := p.cache[key]; prior.TMDBID > 0 && prior.PosterPath != "" &&
		(p.overrides[key] == 0 || p.overrides[key] == prior.TMDBID) {
		if result.TMDBID == 0 && (result.State == "error" || result.State == "not-found") {
			return prior
		}
		if result.TMDBID == prior.TMDBID && result.PosterPath == "" {
			result.PosterPath = prior.PosterPath
		}
        if result.TMDBID==prior.TMDBID&&result.Note!=""&&result.Runtime==0{
         result.Runtime=prior.Runtime;result.Director=prior.Director;result.Genres=prior.Genres;result.Countries=prior.Countries
         if result.VoteCount==0{result.VoteAverage=prior.VoteAverage;result.VoteCount=prior.VoteCount}
         if result.Title==""{result.Title=prior.Title}
        }
	}
	p.putCachedLocked(key, result)
	return result
}

func (p *availabilityServer) putCachedLocked(k string, r providerResult) {
	p.cache[k] = r
	p.cacheDirty = true
	if r.TMDBID > 0 {
		if p.identities == nil {
			p.identities = make(map[string]filmIdentityRecord)
		}
		p.identities[k] = filmIdentityRecord{TMDBID: r.TMDBID, State: "ok", Updated: time.Now().Format(time.RFC3339)}
		p.identityDirty = true
	}
	if p.cacheTimer == nil {
		p.cacheTimer = time.AfterFunc(3*time.Second, p.flushProviderCache)
	}
}

func (p *availabilityServer) flushProviderCache() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flushProviderCacheLocked()
}

func (p *availabilityServer) flushProviderCacheLocked() {
	if p.cacheTimer != nil {
		p.cacheTimer.Stop()
		p.cacheTimer = nil
	}
	if p.identityDirty {
		if err := p.saveIdentityCacheLocked(); err != nil {
			logDiagnostic("Caché de identificaciones pendiente: " + err.Error())
		}
	}
	if !p.cacheDirty || p.cachePath == "" {
		return
	}
	data, err := p.encodeProviderCacheLocked()
	if err == nil {
		err = os.MkdirAll(filepath.Dir(p.cachePath), 0700)
	}
	if err == nil {
		err = replaceLibraryFile(p.cachePath, data)
	}
	if err != nil {
		logDiagnostic("Caché de disponibilidad pendiente: " + err.Error())
		return
	}
	p.cacheDirty = false
}

// A confirmed TMDB ID takes precedence over fuzzy title matches. Keys include
// the release year so similarly named films do not overwrite each other.
func movieCacheKey(title string, year int) string {
	return normaliseTitle(title) + "|" + strconv.Itoa(year)
}
func (p *availabilityServer) overrideFor(title string, year int) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.overrides[movieCacheKey(title, year)]
}
func (p *availabilityServer) saveOverride(title string, year, id int) error {
	key := movieCacheKey(title, year)
	p.mu.Lock()
	defer p.mu.Unlock()
	next := make(map[string]int, len(p.overrides)+1)
	for k, v := range p.overrides {
		next[k] = v
	}
	if id > 0 {
		next[key] = id
	} else {
		delete(next, key)
	}
	bytes, e := json.MarshalIndent(next, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(p.overridesPath), 0700); e != nil {
		return e
	}
	temp := p.overridesPath + ".tmp"
	if e = os.WriteFile(temp, bytes, 0600); e != nil {
		return e
	}
	if e = os.Rename(temp, p.overridesPath); e != nil {
		if e2 := os.WriteFile(p.overridesPath, bytes, 0600); e2 != nil {
			return e2
		}
		_ = os.Remove(temp)
	}
	p.overrides = next
	p.bumpOverrideRevisionLocked(key)
	if id > 0 {
		p.identities[key] = filmIdentityRecord{TMDBID: id, State: "ok", Updated: time.Now().Format(time.RFC3339)}
	} else {
		delete(p.identities, key)
	}
	p.saveIdentityCacheLocked()
	delete(p.cache, key)
	// Invalidating on disk too prevents a stale result if the app exits before
	// the next lookup writes a fresh entry.
	p.cacheDirty = true
	p.flushProviderCacheLocked()
	return nil
}

// Search exact title first, then conservative variants used by distributors.
// E.g. "Vértigo. De entre los muertos" -> "Vértigo" -> "Vertigo".
func titleSearchVariants(title string) []string {
	seen := map[string]bool{}
	var variants []string
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" || len([]rune(s)) < 4 {
			return
		}
		k := strings.ToLower(s)
		if !seen[k] {
			seen[k] = true
			variants = append(variants, s)
		}
	}
	add(title)
	for _, sep := range []string{". ", ": ", " — ", " - "} {
		if i := strings.Index(title, sep); i >= 0 {
			add(title[:i])
		}
	}
	if i := strings.Index(title, " ("); i >= 0 {
		add(title[:i])
	}
	// Accent-free aliases help international TMDB search indexing.
	existing := append([]string(nil), variants...)
	for _, v := range existing {
		ascii := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U").Replace(v)
		add(ascii)
	}
	if len(variants) > 7 {
		variants = variants[:7]
	}
	return variants
}

func newAvailabilityError(s string) providerResult { return providerResult{State: "error", Note: s} }
func requestTMDB(ctx context.Context, path string, params url.Values, key string, v any) (resultErr error) {
	started := time.Now()
	release,err:=acquireTMDBRequest(ctx,tmdbClient)
	if err!=nil{return err}
	defer release()
	defer func() { if resultErr != nil { diagnosticEvent("WARN", "tmdb", fmt.Sprintf("consulta fallida cause=%s ms=%d", diagnosticErrorKind(resultErr), time.Since(started).Milliseconds())) } }()
	if !strings.HasPrefix(key, "eyJ") && !strings.Contains(key, ".") {
		params.Set("api_key", key)
	}
	u := "https://api.themoviedb.org/3" + path + "?" + params.Encode()
	req, e := http.NewRequestWithContext(ctx, "GET", u, nil)
	if e != nil {
		return e
	}
	req.Header.Set("Accept", "application/json")
	if strings.HasPrefix(key, "eyJ") || strings.Contains(key, ".") {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	res, e := tmdbClient.Do(req)
	if e != nil {
		return fmt.Errorf("Conexión TMDB no disponible: %s", diagnosticErrorKind(e))
	}
	defer res.Body.Close()
	if res.StatusCode != 200 { diagnosticEvent("WARN", "tmdb", fmt.Sprintf("TMDB status=%d", res.StatusCode)) }
	if res.StatusCode == 401 || res.StatusCode == 403 {
		return errors.New("Clave TMDB incorrecta o sin autorización")
	}
	if res.StatusCode == 429 {
		deferTMDBRequests(tmdbClient,res.Header.Get("Retry-After"))
		return errors.New("Límite temporal de TMDB. Prueba a actualizar más tarde")
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("TMDB devolvió HTTP %d", res.StatusCode)
	}
	dec := json.NewDecoder(io.LimitReader(res.Body, 2*1024*1024))
	return dec.Decode(v)
}
func normaliseTitle(s string) string {
	s = strings.ToLower(s)
	repl := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n", "à", "a", "è", "e", "ò", "o", "ï", "i", "â", "a", "ê", "e", "ô", "o", "ç", "c", "æ", "ae", "œ", "oe")
	s = repl.Replace(s)
	return strings.Join(strings.FieldsFunc(s, func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsDigit(c) }), " ")
}
func similarTitle(a, b string) float64 {
	a = normaliseTitle(a)
	b = normaliseTitle(b)
	if a == b && a != "" {
		return 1
	}
	if a == "" || b == "" {
		return 0
	}
	if len(a) > 12 && strings.Contains(b, a) || len(b) > 12 && strings.Contains(a, b) {
		return .83
	}
	aw := strings.Fields(a)
	bw := strings.Fields(b)
	sa := map[string]bool{}
	sb := map[string]bool{}
	for _, x := range aw {
		sa[x] = true
	}
	for _, x := range bw {
		sb[x] = true
	}
	match := 0
	for x := range sa {
		if sb[x] {
			match++
		}
	}
	denom := len(sa) + len(sb) - match
	if denom == 0 {
		return 0
	}
	return float64(match) / float64(denom)
}
func yearFromDate(s string) int {
	if len(s) < 4 {
		return 0
	}
	n, _ := strconv.Atoi(s[:4])
	return n
}
func chooseMovie(query string, year int, results []struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Original string `json:"original_title"`
	Released string `json:"release_date"`
}) (int, string) {
	bestID := 0
	bestTitle := ""
	bestScore := 0.0
	for _, m := range results {
		score := similarTitle(query, m.Title)
		if alt := similarTitle(query, m.Original); alt > score {
			score = alt
		}
		y := yearFromDate(m.Released)
		if year > 0 && y > 0 {
			d := year - y
			if d < 0 {
				d = -d
			}
			if d > 2 {
				continue
			}
			if d == 0 {
				score += .1
			} else if d == 1 {
				score += .03
			} else {
				score -= .12
			}
		}
		if score > bestScore {
			bestScore = score
			bestID = m.ID
			bestTitle = m.Title
		}
	}
	if bestScore < .8 {
		return 0, ""
	} // Evita títulos parecidos, pero incorrectos
	return bestID, bestTitle
}
func names(a []struct {
	Name string `json:"provider_name"`
}) []string {
	ret := make([]string, 0, len(a))
	seen := map[string]bool{}
	for _, x := range a {
		n := strings.TrimSpace(x.Name)
		if n != "" && !seen[n] {
			ret = append(ret, n)
			seen[n] = true
		}
	}
	return ret
}
func (p *availabilityServer) fetchFilm(ctx context.Context, title string, year int, force bool) providerResult {
	return p.fetchFilmWithID(ctx, title, year, force, 0)
}

func providerCacheExpired(r providerResult, now time.Time) bool {
 at,err:=time.Parse(time.RFC3339,r.Date)
 if err!=nil{return true}
 ttl:=24*time.Hour
 if r.State=="error"{ttl=15*time.Minute}
 return !now.Before(at.Add(ttl))
}
func (p *availabilityServer) fetchFilmWithID(ctx context.Context, title string, year int, force bool, preferredID int) providerResult {
	cacheKey := movieCacheKey(title, year)
	p.mu.Lock()
	id, revision := p.overrides[cacheKey], p.overrideRevision[cacheKey]
	p.mu.Unlock()
	if id > 0 {
		preferredID = id
	}
	finish := func(result providerResult) providerResult { return p.finishAvailability(cacheKey, revision, result) }
	if c, ok := p.getCached(cacheKey, force); ok && (preferredID == 0 || c.State == "error" || c.TMDBID == preferredID) {
        if !providerCacheExpired(c,time.Now()) || p.getKey()=="" {return c}
        if c.TMDBID>0 {preferredID=c.TMDBID}
	}
	apiKey := p.getKey()
	if apiKey == "" {
		return newAvailabilityError("Configura una clave de TMDB")
	}
	type movie struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Original    string  `json:"original_title"`
		Released    string  `json:"release_date"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		PosterPath  string  `json:"poster_path"`
	}
	var found struct {
		Results []movie `json:"results"`
	}
	matched := ""
	if id == 0 {
		id = preferredID
	}
	if id == 0 {
		// First try the full title, then its likely shorter/original aliases.
		for _, query := range titleSearchVariants(title) {
			var search struct {
				Results []movie `json:"results"`
			}
			params := url.Values{"language": {"es-ES"}, "query": {query}, "include_adult": {"false"}}
			if err := requestTMDB(ctx, "/search/movie", params, apiKey, &search); err != nil {
				r := newAvailabilityError(err.Error())
				r.Date = time.Now().Format(time.RFC3339)
				r.MetadataVersion = 3
				return finish(r)
			}
			choices := make([]struct {
				ID       int    `json:"id"`
				Title    string `json:"title"`
				Original string `json:"original_title"`
				Released string `json:"release_date"`
			}, 0, len(search.Results))
			for _, m := range search.Results {
				choices = append(choices, struct {
					ID       int    `json:"id"`
					Title    string `json:"title"`
					Original string `json:"original_title"`
					Released string `json:"release_date"`
				}{m.ID, m.Title, m.Original, m.Released})
			}
			candidateID, candidateName := chooseMovie(query, year, choices)
			if candidateID != 0 {
				id, matched, found.Results = candidateID, candidateName, search.Results
				break
			}
		}
	}
	if id == 0 {
		r := providerResult{State: "not-found", Date: time.Now().Format(time.RFC3339), MetadataVersion: 3}
		return finish(r)
	}
	type off struct {
		Name string `json:"provider_name"`
	}
	var all struct {
		Results map[string]struct {
			Link     string `json:"link"`
			Flatrate []off  `json:"flatrate"`
			Free     []off  `json:"free"`
			Ads      []off  `json:"ads"`
			Rent     []off  `json:"rent"`
			Buy      []off  `json:"buy"`
		} `json:"results"`
	}
	var details struct {
		Title       string  `json:"title"`
		Released    string  `json:"release_date"`
		Runtime     int     `json:"runtime"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		PosterPath  string  `json:"poster_path"`
		Genres      []struct {
			Name string `json:"name"`
		} `json:"genres"`
		ProductionCountries []struct {
			Name string `json:"name"`
			ISO  string `json:"iso_3166_1"`
		} `json:"production_countries"`
		Credits struct {
			Crew []struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			} `json:"crew"`
		} `json:"credits"`
	}
	r := providerResult{State: "ok", Title: matched, Date: time.Now().Format(time.RFC3339), MetadataVersion: 3, TMDBID: id}
	for _, film := range found.Results {
		if film.ID == id {
			r.VoteAverage = film.VoteAverage
			r.VoteCount = film.VoteCount
			r.PosterPath = film.PosterPath
			break
		}
	}
	if err := requestTMDB(ctx, "/movie/"+strconv.Itoa(id), url.Values{"language": {"es-ES"}, "append_to_response": {"credits"}}, apiKey, &details); err == nil {
		if details.Title != "" {
			r.Title = details.Title
		}
		if details.VoteCount > 0 {
			r.VoteAverage = details.VoteAverage
			r.VoteCount = details.VoteCount
		}
		if details.PosterPath != "" {
			r.PosterPath = details.PosterPath
		}
		r.Runtime = details.Runtime
		for _, g := range details.Genres {
			if g.Name != "" {
				r.Genres = append(r.Genres, g.Name)
			}
		}
		for _, c := range details.ProductionCountries {
			if c.Name != "" {
				r.Countries = append(r.Countries, c.Name)
			}
		}
		for _, c := range details.Credits.Crew {
			if c.Job == "Director" && c.Name != "" {
				if r.Director != "" {
					r.Director += " / "
				}
				r.Director += c.Name
			}
		}
	} else {
		// Conservar la ficha parcial y permitir actualización manual posterior.
		r.Note = "Ficha parcial; pulsa Actualizar para reintentar los datos incompletos"
	}
	if !validPosterPath(r.PosterPath) {
		r.PosterPath = ""
	}
	if e := requestTMDB(ctx, "/movie/"+strconv.Itoa(id)+"/watch/providers", url.Values{}, apiKey, &all); e != nil {
		r.State = "error"
		r.Note = e.Error()
		return finish(r)
	}
	r.Date = time.Now().Format(time.RFC3339)
	es, ok := all.Results["ES"]
	r.FilmURL = es.Link
	if !ok {
		r.State = "no-data"
	} else {
		convert := func(arr []off) []string {
			out := make([]string, 0, len(arr))
			seen := map[string]bool{}
			for _, x := range arr {
				if x.Name != "" && !seen[x.Name] {
					out = append(out, x.Name)
					seen[x.Name] = true
				}
			}
			return out
		}
		r.Subscriptions = convert(es.Flatrate)
		r.Free = convert(es.Free)
		r.Ads = convert(es.Ads)
		r.Rent = convert(es.Rent)
		r.Buy = convert(es.Buy)
	}
	return finish(r)
}
func (a *runningApp) servePlatformHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	switch r.URL.Path {
	case "/api/platforms/overrides":
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"overrides": a.platforms.manualOverrides()})
	case "/api/platforms/overrides/import":
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var payload struct {
			Overrides map[string]int `json:"overrides"`
		}
		if e := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&payload); e != nil {
			http.Error(w, "Identificaciones inválidas", 400)
			return
		}
		if e := a.platforms.importManualOverrides(payload.Overrides); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"saved": true})
	case "/api/platforms/cache":
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"entries": a.platforms.cachedEntries()})
	case "/api/platforms/status":
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"configured": a.platforms.configured()})
	case "/api/platforms/config":
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var body struct {
			Key string `json:"key"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body) != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		body.Key = strings.TrimSpace(body.Key)
		if len(body.Key) < 12 || len(body.Key) > 512 {
			http.Error(w, "Clave no válida", 400)
			return
		}
		// Verifica credenciales antes de conservarlas en disco.
		ctx, cancel := context.WithTimeout(r.Context(), 18*time.Second)
		defer cancel()
		var result map[string]any
		if err := requestTMDB(ctx, "/configuration", url.Values{}, body.Key, &result); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := a.platforms.saveKey(body.Key); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"configured": true})
	case "/api/platforms/match":
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var body struct {
			Title              string `json:"title"`
			Year               int    `json:"year"`
			ID                 int    `json:"id"`
			AcceptYearMismatch bool   `json:"accept_year_mismatch"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body); err != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		body.Title = strings.TrimSpace(body.Title)
		if body.Title == "" || len(body.Title) > 250 || body.Year < 0 || body.Year > 2200 || body.ID < 0 || body.ID > 999999999 {
			http.Error(w, "Identificación inválida", 400)
			return
		}
		if body.ID > 0 {
			var film struct {
				Title    string `json:"title"`
				Released string `json:"release_date"`
			}
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()
			err := requestTMDB(ctx, "/movie/"+strconv.Itoa(body.ID), url.Values{"language": {"es-ES"}}, a.platforms.getKey(), &film)
			if err != nil {
				http.Error(w, "No se ha validado el ID de TMDB: "+err.Error(), 400)
				return
			}
			if film.Title == "" {
				http.Error(w, "El ID no corresponde a una película", 400)
				return
			}
			if year := yearFromDate(film.Released); body.Year > 0 && year > 0 && !body.AcceptYearMismatch {
				delta := body.Year - year
				if delta < -2 || delta > 2 {
					_ = json.NewEncoder(w).Encode(map[string]any{"saved": false, "year_mismatch": year})
					return
				}
			}
		}
		if err := a.platforms.saveOverride(body.Title, body.Year, body.ID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"saved": true})
	case "/api/platforms/lookup":
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", 405)
			return
		}
		var body struct {
			Title  string `json:"title"`
			Year   int    `json:"year"`
			Force  bool   `json:"force"`
			TMDBID int    `json:"tmdb_id"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body) != nil {
			http.Error(w, "JSON inválido", 400)
			return
		}
		body.Title = strings.TrimSpace(body.Title)
		if body.Title == "" || len(body.Title) > 250 || body.Year < 0 || body.Year > 2200 || body.TMDBID < 0 || body.TMDBID > 999999999 {
			http.Error(w, "Película no válida", 400)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
		defer cancel()
		result := a.platforms.fetchFilmWithID(ctx, body.Title, body.Year, body.Force, body.TMDBID)
		_ = json.NewEncoder(w).Encode(result)
	default:
		http.NotFound(w, r)
	}
}

// Images come exclusively from TMDB and are cached locally for offline display.
var safePoster = regexp.MustCompile(`^/[A-Za-z0-9_-]{5,80}\.(?:jpg|jpeg|png|webp)$`)

func validPosterPath(path string) bool { return safePoster.MatchString(path) }
func (a *runningApp) servePosterHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !a.auth(r) {
		http.Error(w, "No autorizado", 403)
		return
	}
	path := r.URL.Query().Get("path")
	if !validPosterPath(path) {
		http.Error(w, "Cartel no válido", 400)
		return
	}
	file := filepath.Join(a.platforms.posterDir, strings.TrimPrefix(path, "/"))
	if b, e := os.ReadFile(file); e == nil {
		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
		w.Header().Set("Cache-Control", "private, max-age=604800")
		_, _ = w.Write(b)
		return
	}
	if desktopFeatures {
		if data := a.discoveryImages.get(discoveryImageDir(a.libraryPath), filepath.Base(path)); len(data) > 0 {
			w.Header().Set("Content-Type", discoveryImageType(data))
			w.Header().Set("Cache-Control", "private, max-age=604800")
			_, _ = w.Write(data)
			return
		}
	}
	// Copia bajo demanda el cartel de Roaming si existe; otra versión puede
	// seguir usando ese directorio, por eso no se borra su copia.
	if a.platforms.legacyPosterDir != "" {
		old := filepath.Join(a.platforms.legacyPosterDir, strings.TrimPrefix(path, "/"))
		if b, err := os.ReadFile(old); err == nil && len(b) > 100 && len(b) <= 1<<20 {
			_ = os.MkdirAll(a.platforms.posterDir, 0700)
			_ = os.WriteFile(file, b, 0600)
			if protectedLibraryEnabled {
				_, _ = prunePosterCache(a.platforms.posterDir, appPosterMaxFiles, appPosterMaxBytes, appPosterMaxAge)
			}
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
			w.Header().Set("Cache-Control", "private, max-age=604800")
			_, _ = w.Write(b)
			return
		}
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	req, e := http.NewRequestWithContext(ctx, "GET", "https://image.tmdb.org/t/p/w92"+path, nil)
	if e != nil {
		http.Error(w, "Cartel no disponible", 502)
		return
	}
	res, e := tmdbClient.Do(req)
	if e != nil {
		http.Error(w, "Sin conexión", 502)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		http.Error(w, "Cartel no disponible", 502)
		return
	}
	body, e := io.ReadAll(io.LimitReader(res.Body, 1024*1024+1))
	if e != nil || len(body) > 1024*1024 || len(body) < 100 {
		http.Error(w, "Imagen no válida", 502)
		return
	}
	_ = os.MkdirAll(a.platforms.posterDir, 0700)
	_ = os.WriteFile(file, body, 0600)
	if protectedLibraryEnabled {
		_, _ = prunePosterCache(a.platforms.posterDir, appPosterMaxFiles, appPosterMaxBytes, appPosterMaxAge)
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(file)))
	_, _ = w.Write(body)
}
