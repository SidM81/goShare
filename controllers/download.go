package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/SidM81/goShare/config"
	"github.com/SidM81/goShare/models"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// This function will handle file downloads
	// It will retrieve the file from MinIO and send it to the client
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "Missing file name", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	var file models.File
	if err := db.First(&file, "name = ?", fileName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "DB error", http.StatusInternalServerError)
		}
		return
	}

	// Fetch and sort chunks
	var chunks []models.Chunk
	if err := db.Where("file_id = ?", file.ID).Find(&chunks).Error; err != nil {
		http.Error(w, "Error fetching chunks", http.StatusInternalServerError)
		return
	}
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Index < chunks[j].Index
	})

	minioClient := config.GetMinioClient()
	bucket := config.MustHaveConfig().MINIO_BUCKET_NAME

	var fullData bytes.Buffer
	for _, chunk := range chunks {
		objectName := fmt.Sprintf("%s_%d", file.ID, chunk.Index)

		obj, err := minioClient.GetObject(r.Context(), bucket, objectName, minio.GetObjectOptions{})
		if err != nil {
			http.Error(w, "Error retrieving chunk from storage", http.StatusInternalServerError)
			return
		}
		defer obj.Close()

		if _, err := io.Copy(&fullData, obj); err != nil {
			http.Error(w, "Error reading chunk", http.StatusInternalServerError)
			return
		}
	}

	// Send file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fullData.Bytes())
}
