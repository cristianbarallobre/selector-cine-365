//go:build windows || darwin

package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDesktopDirectorCreditsAndFilmography(t *testing.T) {
	oldClient := tmdbClient
	defer func() { tmdbClient = oldClient }()
	tmdbClient = &http.Client{Transport: roundTripMock(func(r *http.Request) (*http.Response, error) {
		var body string
		switch r.URL.Path {
		case "/3/movie/42":
			body = `{"id":42,"title":"Película","credits":{"crew":[
				{"id":101,"name":"Directora","job":"Director","profile_path":"/director_a.jpg"},
				{"id":101,"name":"Directora","job":"Director","profile_path":"/director_a.jpg"},
				{"id":102,"name":"Codirector","job":"Director","profile_path":"https://example.invalid/retrato.jpg"},
				{"id":103,"name":"Guionista","job":"Writer"}]}}`
		case "/3/movie/42/videos":
			body = `{"results":[]}`
		case "/3/movie/42/watch/providers":
			body = `{"results":{}}`
		case "/3/person/101":
			body = `{"id":101,"name":"Directora","biography":"Biografía de prueba","profile_path":"/director_a.jpg",
			 "movie_credits":{"cast":[{"id":77,"title":"Como intérprete"}],"crew":[
			 {"id":42,"title":"Película dirigida","job":"Director"},{"id":42,"title":"Copia","job":"Director"},
			 {"id":88,"title":"Solo producción","job":"Producer"}]}}`
		default:
			t.Errorf("unexpected TMDB request: %s", r.URL.Path)
			return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("{}"))}, nil
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	detail, err := labFreshDetail(context.Background(), 42, "fixture", labDetailRecord{})
	if err != nil {
		t.Fatal(err)
	}
	directors, ok := detail.Data["directors"].([]labPerson)
	if !ok || len(directors) != 2 || directors[0].ID != 101 || directors[0].ProfilePath != "/director_a.jpg" || directors[1].ProfilePath != "" {
		t.Fatalf("unexpected director credits: %+v", detail.Data["directors"])
	}
	if detail.Data["director"] != "Directora / Codirector" {
		t.Fatal("the legacy director label was lost")
	}

	app := &runningApp{platforms: &availabilityServer{apiKey: "fixture"}}
	for _, tc := range []struct {
		query string
		id    int
	}{{"&role=director", 42}, {"", 77}} {
		response := httptest.NewRecorder()
		app.labPersonDetail(response, httptest.NewRequest(http.MethodGet, "/api/discover/person?id=101"+tc.query, nil))
		var payload struct {
			Movies []labPersonFilm `json:"movies"`
		}
		if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &payload) != nil || len(payload.Movies) != 1 || payload.Movies[0].ID != tc.id {
			t.Fatalf("unexpected filmography for %q: %d %s", tc.query, response.Code, response.Body.String())
		}
	}
}
