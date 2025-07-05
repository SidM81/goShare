package routes

import "net/http"

func SetupRouter() *http.ServeMux {
	r := http.NewServeMux()

	// File upload route
	r.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("File upload functionality is not implemented yet."))
	})

	// File download route
	r.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("File download functionality is not implemented yet."))
	})

	// File share route
	r.HandleFunc("/share/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("File share functionality is not implemented yet."))
	})

	return r
}
