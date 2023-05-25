package upload

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type TelegraPhResponse struct {
	Src string `json:"src"`
}

func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "error", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "error", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "error", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://telegra.ph/upload", body)
	if err != nil {
		return "error", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "error", err
	}
	defer resp.Body.Close()

	var telegraPhResponse []TelegraPhResponse
	err = json.NewDecoder(resp.Body).Decode(&telegraPhResponse)
	if err != nil {
		return "error", err
	}

	return telegraPhResponse[0].Src, nil
}
