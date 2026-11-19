package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripMock func(*http.Request) (*http.Response, error)

func (f roundTripMock) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestPlatformMatchingAndProviderKinds(t *testing.T) {
	old := tmdbClient
	defer func() { tmdbClient = old }()
	calls := 0
	tmdbClient = &http.Client{Timeout: 5 * time.Second, Transport: roundTripMock(func(r *http.Request) (*http.Response, error) {
		calls++
		body := ""
		if strings.HasSuffix(r.URL.Path, "/search/movie") {
			body = `{"results":[{"id":10,"title":"Blade Runner","original_title":"Blade Runner","release_date":"1982-06-25","vote_average":8.1,"vote_count":1800,"poster_path":"/qwerty12345.jpg"},{"id":11,"title":"Blade Runner 2049","original_title":"Blade Runner 2049","release_date":"2017-10-04"}]}`
		} else if strings.HasSuffix(r.URL.Path, "/movie/10") {
			body = `{"runtime":117,"vote_average":8.1,"vote_count":1800,"poster_path":"/qwerty12345.jpg","genres":[{"name":"Ciencia ficción"}],"production_countries":[{"name":"Estados Unidos"}],"credits":{"crew":[{"name":"Ridley Scott","job":"Director"}]}}`
		} else if strings.HasSuffix(r.URL.Path, "/movie/10/watch/providers") {
			body = `{"results":{"ES":{"link":"https://www.themoviedb.org/movie/10/watch?locale=ES","flatrate":[{"provider_name":"Filmin"}],"rent":[{"provider_name":"Apple TV"}]}}}`
		} else {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
	p := newAvailabilityServer(t.TempDir())
	if err := p.saveKey(strings.Repeat("a", 32)); err != nil {
		t.Fatal(err)
	}
	v := p.fetchFilm(context.Background(), "Blade Runner", 1982, false)
	if v.State != "ok" || len(v.Subscriptions) != 1 || v.Subscriptions[0] != "Filmin" || len(v.Rent) != 1 {
		t.Fatalf("unexpected result: %+v", v)
	}
	c := p.fetchFilm(context.Background(), "Blade Runner", 1982, false)
	if c.State != "ok" || calls != 3 {
		t.Fatalf("cache missed: %d %+v", calls, c)
	}
	if v.VoteAverage != 8.1 || v.Runtime != 117 || v.Director != "Ridley Scott" || len(v.Genres) != 1 || v.PosterPath == "" {
		t.Fatalf("missing movie metadata: %+v", v)
	}
	p2 := newAvailabilityServer(t.TempDir())
	if p2.configured() {
		t.Fatal("config not isolated")
	}
}
func TestStrictMatching(t *testing.T) {
	type candidate = struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Original string `json:"original_title"`
		Released string `json:"release_date"`
	}
	id, _ := chooseMovie("Vertigo", 1958, []candidate{{17, "Vértigo", "Vertigo", "1958-03-01"}, {12, "Vertigo", "Vertigo", "2015-01-01"}})
	if id != 17 {
		t.Fatalf("wanted 1958 film, got %d", id)
	}
	id, _ = chooseMovie("Atrapado en el tiempo", 1993, []candidate{{22, "Atrapado en el tiempo II", "Atrapado en el tiempo II", "2003-01-01"}})
	if id != 0 {
		t.Fatalf("must not guess from wrong year: %d", id)
	}
}
