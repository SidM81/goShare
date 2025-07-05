# GoShare — Chunked File Sharing System in Go

GoShare is a minimalist backend system written in **Go** that allows users to upload large files, automatically chunks and hashes them, stores them on **MinIO**, and allows secure downloading by filename. File metadata is stored in a **SQLite** database.

## Features

- File upload with automatic:
  - Chunking (4MB per chunk)
  - SHA-256 hashing per chunk and full file
- Chunk storage on MinIO under `<fileID>/<chunkIndex>`
- Metadata stored in SQLite (`File` and `Chunk` models)
- Download file by name — chunks reassembled server-side
- Shareable download URL generated after upload

## Tech Stack

- **Go** (net/http, GORM)
- **SQLite** (via GORM)
- **MinIO** (object storage)
- **WSL or Linux** for development
- **cURL** or HTTP clients to interact

## Project Structure

```
goShare/
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── controllers/
│   ├── upload.go
│   ├── download.go
│   └── share.go
├── models/
│   ├── file.go
│   ├── share.go
│   └── chunk.go
├── routes/
│   └── router.go
├── go.mod
├── .env
└── README.md
```
## Running the Project

### 1. Clone and Setup

```
git clone https://github.com/SidM81/goShare.git
cd goShare
go mod tidy
```

### 2. Set up `.env`

Create a `.env` file at root with your MinIO + DB config (see config file).

### 3. Start MinIO Server

```
docker run -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  quay.io/minio/minio server /data --console-address ":9001"
```

### 4. Run Go Server

```
go run cmd/server/main.go
```

## API Usage

### Upload File

```
curl -X POST http://localhost:8080/upload/ \
  -F "file=@/path/to/your/file.pdf" \
  -F "name=your_custom_name.pdf"
```

**Response:**

```
{
  "message": "Upload complete",
  "file_id": "uuid-here",
  "download_url": "http://localhost:8080/download?name=your_custom_name.pdf"
}
```

### Download File

```
curl -O -J "http://localhost:8080/download?name=your_custom_name.pdf"
```

The `-O -J` ensures filename from headers is used.

## Notes

- Files are chunked at 4MB size for better handling of large uploads.
- Make sure your MinIO bucket is created (`goshare`) before uploading.
- Use MinIO Console ([http://localhost:9001](http://localhost:9001)) to inspect files if needed.

## Author

**Siddharth Mishra**\
Backend Developer | Competitive Programmer\
[GitHub](https://github.com/SidM81)

## License

This project is licensed under the MIT License.
