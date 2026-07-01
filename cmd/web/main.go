package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.AppName}}</title>
  <style>
    body { font-family: system-ui, sans-serif; margin: 3rem; color: #172026; background: #f6f8fb; }
    main { max-width: 760px; background: white; border: 1px solid #d9e1ea; border-radius: 8px; padding: 2rem; }
    dt { font-weight: 700; margin-top: 1rem; }
    dd { margin: .25rem 0 0; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  </style>
</head>
<body>
<main>
  <h1>{{.AppName}}</h1>
  <dl>
    <dt>Preview PR</dt><dd>{{.PreviewPR}}</dd>
    <dt>Version</dt><dd>{{.Version}}</dd>
    <dt>Hostname</dt><dd>{{.Hostname}}</dd>
    <dt>Shared secret</dt><dd>{{.SharedSecret}}</dd>
    <dt>Generated secret</dt><dd>{{.GeneratedSecret}}</dd>
    <dt>Started</dt><dd>{{.Started}}</dd>
  </dl>
</main>
</body>
</html>`))

type data struct {
	AppName         string
	PreviewPR       string
	Version         string
	Hostname        string
	SharedSecret    string
	GeneratedSecret string
	Started         string
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	started := time.Now().UTC().Format(time.RFC3339)
	hostname, _ := os.Hostname()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		_ = page.Execute(w, data{
			AppName:         env("APP_NAME", "web-apps"),
			PreviewPR:       env("PREVIEW_PR", "none"),
			Version:         env("APP_VERSION", "dev"),
			Hostname:        hostname,
			SharedSecret:    env("SHARED_SECRET_MESSAGE", "missing"),
			GeneratedSecret: env("GENERATED_SECRET_MESSAGE", "missing"),
			Started:         started,
		})
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; version=0.0.4")
		_, _ = fmt.Fprintf(w, "web_apps_info{preview_pr=%q,version=%q} 1\n", env("PREVIEW_PR", "none"), env("APP_VERSION", "dev"))
	})

	addr := env("HTTP_ADDR", ":8080")
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

