package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SidM81/goShare/config"
	"github.com/SidM81/goShare/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// This function will handle file uploads
	// It will receive the file from the request, save it to MinIO, and store metadata in SQLite
	fmt.Println("File Upload Endpoint Hit")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20) // Limit upload size to 32MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("name")
	if fileName == "" {
		fileName = handler.Filename
	}

	fileID := uuid.New()
	minioClient := config.GetMinioClient()
	bucket := config.MustHaveConfig().MINIO_BUCKET_NAME
	db := config.GetDB()

	hasher := sha256.New()
	chunkIndex := 0
	totalSize := int64(0)
	chunkBuf := make([]byte, 4*1024*1024) // 4MB buffer size

	for {
		n, err := file.Read(chunkBuf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		chunkData := chunkBuf[:n]
		hasher.Write(chunkData)
		chunkHash := sha256.Sum256(chunkData)
		chunkID := uuid.New()
		chunkKey := fmt.Sprintf("%s/%d", fileID.String(), chunkIndex)

		_, err = minioClient.PutObject(r.Context(), bucket, chunkKey, bytes.NewReader(chunkData), int64(n), minio.PutObjectOptions{})
		if err != nil {
			http.Error(w, "Failed to upload chunk to MinIO", http.StatusInternalServerError)
			return
		}

		chunk := models.Chunk{
			ID:       chunkID,
			FileID:   fileID,
			Index:    chunkIndex,
			Hash:     hex.EncodeToString(chunkHash[:]),
			Uploaded: true,
		}
		db.Create(&chunk)

		totalSize += int64(n)
		chunkIndex++
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))
	uploadedFile := models.File{
		ID:        fileID,
		Name:      fileName,
		Size:      totalSize,
		NumChunks: chunkIndex,
		Hash:      fileHash,
		Status:    "complete",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&uploadedFile)

	// Include download link in the response
	downloadURL := fmt.Sprintf("http://localhost:8080/download?name=%s", fileName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"message": "Upload complete",
		"file_id": "%s",
		"download_url": "%s"
	}`, fileID.String(), downloadURL)))

}
