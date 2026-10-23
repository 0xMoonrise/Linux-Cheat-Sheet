package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const page = `<!doctype html>
<meta charset=utf-8>
<meta name=viewport content="width=device-width,initial-scale=1">
<link rel=stylesheet href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.8.1/github-markdown.min.css">
<link rel=stylesheet media="(prefers-color-scheme: light)" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/github.min.css">
<link rel=stylesheet media="(prefers-color-scheme: dark)" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/github-dark.min.css">
<style>
  body { margin: 0 }
  .markdown-body { box-sizing: border-box; max-width: 980px; margin: 0 auto; padding: 45px }
  @media (max-width: 767px) { .markdown-body { padding: 15px } }
  @media (prefers-color-scheme: dark) { body { background: #0d1117 } }
</style>
<article id=c class=markdown-body>{{CONTENT}}</article>
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/highlight.min.js"></script>
<script>
const hl = () => document.querySelectorAll("pre code").forEach(el => hljs.highlightElement(el));
hl();
new EventSource("/events").onmessage = () =>
  fetch("/render").then(r => r.text()).then(h => { c.innerHTML = h; hl() });
</script>`

var (
	mu      sync.Mutex
	clients = map[chan struct{}]bool{}
)

func render(path string) string {
	src, _ := os.ReadFile(path)
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	md.Convert(src, &buf)
	return buf.String()
}

func main() {
	path := os.Args[1]

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w.Add(filepath.Dir(path))

	go func() {
		for ev := range w.Events {
			if filepath.Base(ev.Name) != filepath.Base(path) {
				continue
			}
			mu.Lock()
			for c := range clients {
				select {
				case c <- struct{}{}:
				default:
				}
			}
			mu.Unlock()
		}
	}()

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(rw, strings.Replace(page, "{{CONTENT}}", render(path), 1))
	})

	http.HandleFunc("/render", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, render(path))
	})

	http.HandleFunc("/events", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/event-stream")
		c := make(chan struct{}, 1)
		mu.Lock()
		clients[c] = true
		mu.Unlock()
		defer func() { mu.Lock(); delete(clients, c); mu.Unlock() }()

		f := rw.(http.Flusher)
		for {
			select {
			case <-c:
				fmt.Fprint(rw, "data: x\n\n")
				f.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	log.Println("http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
