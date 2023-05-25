# Gin Telegraph File Upload

This is a Go application built with Gin to upload media files to telegra.ph and generate shareable URLs. The advantage is that all files are stored on the Telegraph server and do not take up space on your server.

## Prerequisites

- Go
- Redis

## Installation

1. Clone the repository:

```shell
git clone https://github.com/stepashka20/gin-telegraph-media-upload.git
```

2. Install dependencies:

```shell
go get
```

3. Copy the `.env.example` file to `.env` and fill in the required values:

```shell
cp .env.example .env
```

4. Run the application:

```shell
go run main.go
```

## Usage

The application exposes the following endpoints:

POST /upload: Upload a file to telegra.ph and get a shareable URL. 5MB file size limit (telegra.ph limit)
GET /:key: Retrieve the file using the generated key.

## License
This project is licensed under the MIT License.
