package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"gin-telegraph-media-upload/upload"
	"path/filepath"
)

var ctx = context.Background()

func rndStr(l int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func main() {
	rand.Seed(time.Now().UnixNano())

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(limits.RequestSizeLimiter(5 << 20)) // telegra.ph limit 5 MB

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		possibleTypes := []string{"image/jpeg", "image/png", "image/gif", "video/mp4", "video/m4v"} // telegra.ph support only this types
		isValid := false
		for _, t := range possibleTypes {
			if t == file.Header.Get("Content-Type") {
				isValid = true
				break
			}
		}
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
			return
		}

		result, err := upload.UploadFile(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fileExt := filepath.Ext(file.Filename)
		key := rndStr(10) + fileExt
		rdb.Set(ctx, key, result, 0)
		c.JSON(http.StatusOK, gin.H{"result": key})
	})

	router.GET("/:key", func(c *gin.Context) {
		value, err := rdb.Get(ctx, c.Param("key")).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "file not found"})
			return
		}
		proxyURL := "https://telegra.ph" + value
		fmt.Println(proxyURL)
		response, err := http.Get(proxyURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer response.Body.Close()

		contentType := response.Header.Get("Content-Type")
		c.DataFromReader(http.StatusOK, response.ContentLength, contentType, response.Body, nil)
	})

	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}
