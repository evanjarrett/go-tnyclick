package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"io"

	"github.com/gorilla/mux"
)

const (
	mediaDir = "./media"
)

func main() {
	// Create the 'media' directory if it doesn't exist
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		err = os.Mkdir(mediaDir, 0755)
		if err != nil {
			log.Fatalf("Error creating 'media' directory: %s", err)
		}
	}

	// Get token from the environment
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("AUTH_TOKEN environment variable not set")
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/images", handleImageUpload(authToken)).Methods("POST")
	r.HandleFunc("/i/{filename}", handleImageRequest).Methods("GET")
	http.Handle("/", r)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleImageUpload(authToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Compare the Authorization header with the environment token
		requestToken := r.Header.Get("Authorization")
		if requestToken != "Token "+authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err := r.ParseMultipartForm(50 << 20) // 10 MB file size limit
		if err != nil {
			http.Error(w, "Unable to process the uploaded file", http.StatusInternalServerError)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate a unique filename for the image using a combination of timestamp and random number
		ext := filepath.Ext(handler.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		newFilePath := filepath.Join(mediaDir, newFileName)

		f, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}

		// Construct the public URL for the image
		publicURL := "https://" + r.Host + "/i/" + newFileName

		w.Write([]byte(publicURL))
	}
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join(mediaDir, filename)
	_, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}
i.j5t.io