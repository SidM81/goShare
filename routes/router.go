package routes

import (
	"net/http"

	"github.com/SidM81/goShare/controllers"
)

func SetupRouter() *http.ServeMux {
	r := http.NewServeMux()

	// File upload route
	r.HandleFunc("/upload/", controllers.UploadFileHandler)

	// File download route
	r.HandleFunc("/download/", controllers.DownloadFileHandler)

	// File share route
	r.HandleFunc("/share/", controllers.ShareFileHandler)

	return r
}
