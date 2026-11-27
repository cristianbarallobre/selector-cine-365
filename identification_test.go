package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLateAvailabilityCannotUndoManualOverride(t *testing.T) {
	p := newAvailabilityServer(t.TempDir())
	p.apiKey = "fixture"
	oldClient := tmdbClient
	defer func() { tmdbClient = oldClient }()
	defer p.flushProviderCache()
	started, release := make(chan struct{}), make(chan struct{})
	finished := make(chan providerResult, 1)
	tmdbClient = &http.Client{Transport: roundTripMock(func(r *http.Request) (*http.Response, error) {
		var body string
		switch r.URL.Path {
		case "/3/search/movie":
			body = `{"results":[{"id":426,"title":"Vértigo","release_date":"1958-01-01"}]}`
		case "/3/movie/426", "/3/movie/999":
			body = `{"title":"Vértigo","release_date":"1958-01-01"}`
		case "/3/movie/426/watch/providers":
			close(started)
			<-release
			body = `{"results":{"ES":{}}}`
		case "/3/movie/999/watch/providers":
			body = `{"results":{"ES":{}}}`
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			body = "{}"
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
	go func() { finished <- p.fetchFilm(context.Background(), "Vértigo", 1958, true) }()
	select {
	case <-started:
	case result := <-finished:
		t.Fatal(result)
	}
	err := p.importManualOverrides(map[string]int{"vertigo|1958": 999})
	close(release)
	outdated := <-finished
	if err != nil || outdated.State != "superseded" {
		t.Fatal(outdated, err)
	}
	if result := p.fetchFilmWithID(context.Background(), "Vértigo", 1958, false, 426); result.TMDBID != 999 {
		t.Fatalf("manual ID lost: %+v", result)
	}
}

func TestManualOverrideImportChangesCacheOnlyAfterCommit(t *testing.T) {
	p := newAvailabilityServer(t.TempDir())
	const key = "vertigo|1958"
	p.cache[key] = providerResult{TMDBID: 426, MetadataVersion: 3}
	if err := os.Mkdir(p.overridesPath, 0700); err != nil {
		t.Fatal(err)
	}
	if err := p.importManualOverrides(map[string]int{key: 999}); err == nil {
		t.Fatal("expected blocked destination to reject the import")
	}
	if got, ok := p.getCached(key, false); !ok || got.TMDBID != 426 {
		t.Fatalf("failed import changed the active availability cache: %v", got)
	}
	if got, _ := p.knownIdentity(key); got != 426 || len(p.manualOverrides()) != 0 {
		t.Fatalf("failed import changed the active identity: %d", got)
	}
	if err := os.Remove(p.overridesPath); err != nil {
		t.Fatal(err)
	}
	if err := p.importManualOverrides(map[string]int{key: 999, "invalid": 12}); err == nil {
		t.Fatal("expected invalid entry to reject the batch")
	}
	if got, ok := p.getCached(key, false); !ok || got.TMDBID != 426 {
		t.Fatal("invalid batch changed the active cache")
	}
	if err := p.importManualOverrides(map[string]int{key: 999}); err != nil {
		t.Fatal(err)
	}
	if _, ok := p.getCached(key, false); ok {
		t.Fatal("successful import retained old availability")
	}
	reopened := newAvailabilityServer(filepath.Dir(p.overridesPath))
	if got, _ := reopened.knownIdentity(key); got != 999 {
		t.Fatalf("reopened identity: %d", got)
	}
}

func TestVertigoFallbackByYear(t *testing.T) {
	variants := titleSearchVariants("Vértigo. De entre los muertos")
	if len(variants) < 3 || variants[0] != "Vértigo. De entre los muertos" || variants[1] != "Vértigo" || variants[len(variants)-1] != "Vertigo" {
		t.Fatalf("unexpected search variants: %#v", variants)
	}
	old := tmdbClient
	defer func() { tmdbClient = old }()
	var queries []string
	tmdbClient = &http.Client{Timeout: 5 * time.Second, Transport: roundTripMock(func(r *http.Request) (*http.Response, error) {
		b := ""
		switch {
		case strings.HasSuffix(r.URL.Path, "/search/movie"):
			q := r.URL.Query().Get("query")
			queries = append(queries, q)
			if q == "Vértigo" {
				b = `{"results":[{"id":426,"title":"Vértigo","original_title":"Vertigo","release_date":"1958-05-09","vote_average":8.0,"vote_count":900,"poster_path":"/qwerty12345.jpg"},{"id":999,"title":"Vértigo","original_title":"Vertigo","release_date":"2015-01-01"}]}`
			} else {
				b = `{"results":[]}`
			}
		case strings.HasSuffix(r.URL.Path, "/movie/426"):
			b = `{"title":"Vértigo","release_date":"1958-05-09","runtime":128,"vote_average":8.0,"vote_count":900,"poster_path":"/qwerty12345.jpg"}`
		case strings.HasSuffix(r.URL.Path, "/movie/426/watch/providers"):
			b = `{"results":{"ES":{"link":"https://www.themoviedb.org/movie/426/watch","flatrate":[{"provider_name":"Filmin"}]}}}`
		default:
			t.Fatalf("unexpected request %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(b))}, nil
	})}
	p := newAvailabilityServer(t.TempDir())
	if err := p.saveKey(strings.Repeat("a", 32)); err != nil {
		t.Fatal(err)
	}
	result := p.fetchFilm(context.Background(), "Vértigo. De entre los muertos", 1958, false)
	if result.State != "ok" || result.TMDBID != 426 || result.Title != "Vértigo" || result.VoteCount != 900 || len(result.Subscriptions) != 1 {
		t.Fatalf("unexpected fallback result: %+v", result)
	}
	if len(queries) != 2 || queries[0] != "Vértigo. De entre los muertos" || queries[1] != "Vértigo" {
		t.Fatalf("unexpected queries: %#v", queries)
	}
	_ = p.fetchFilm(context.Background(), "Vértigo. De entre los muertos", 1958, false)
	if len(queries) != 2 {
		t.Fatal("expected cached result on repeat")
	}
}

func TestManualTMDBIDPersistsAndOverridesOldNegativeCache(t *testing.T) {
	folder := t.TempDir()
	p := newAvailabilityServer(folder)
	key := movieCacheKey("La película desconocida", 1958)
	p.putCached(key, providerResult{State: "not-found", Date: time.Now().Format(time.RFC3339), MetadataVersion: 2})
	if _, ok := p.getCached(key, false); ok {
		t.Fatal("older not-found results must be retried")
	}
	if err := p.saveOverride("La película desconocida", 1958, 426); err != nil {
		t.Fatal(err)
	}
	p2 := newAvailabilityServer(folder)
	if p2.overrideFor("La película desconocida", 1958) != 426 {
		t.Fatal("manual id not persisted")
	}
	if _, ok := p2.getCached(key, false); ok {
		t.Fatal("manual correction must invalidate old cache")
	}
	data, err := os.ReadFile(filepath.Join(folder, "identificaciones_tmdb.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vals map[string]int
	if err := json.Unmarshal(data, &vals); err != nil {
		t.Fatal(err)
	}
	if vals[key] != 426 {
		t.Fatalf("unexpected overrides: %+v", vals)
	}
	if err := p2.saveOverride("La película desconocida", 1958, 0); err != nil {
		t.Fatal(err)
	}
	p3 := newAvailabilityServer(folder)
	if p3.overrideFor("La película desconocida", 1958) != 0 {
		t.Fatal("manual id could not be cleared")
	}
}
