package app

import "net/http"

func Router() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	return m
}
