package app

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const defaultLoadDuration = "10s"

func Router() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	m.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		ds := r.URL.Query().Get("time")
		if ds == "" {
			ds = defaultLoadDuration
		}

		d, err := time.ParseDuration(ds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dl := time.Now().Add(d)

		wg := &sync.WaitGroup{}

		for i := 0; i < runtime.GOMAXPROCS(0); i++ {
			wg.Add(1)
			go func() {
				runtime.LockOSThread()
				// endless loop
				for {
					if dl.Before(time.Now()) {
						break
					}
				}
				wg.Done()
			}()
		}

		wg.Wait()

		_, _ = w.Write([]byte("CPU Done"))
	})

	m.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		ds := r.URL.Query().Get("time")
		if ds == "" {
			ds = defaultLoadDuration
		}

		d, err := time.ParseDuration(ds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		<-ctx.Done()

		_, _ = w.Write([]byte("CPU Done"))
	})

	return m
}
