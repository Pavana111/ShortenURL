package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original"`
	ShortURL     string    `json:"shorturl"`
	CreationTime time.Time `json:"time"`
}

var urldb = make(map[string]URL)

func generateurl(OriginalURL string) string {
	hasher := md5.New()
	fmt.Println("hasher is", hasher)
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	fmt.Println("data is", data)
	hash := hex.EncodeToString(data)
	return hash[:8]

}

func cretaeurl(originalURL string) string {
	shorturl := generateurl(originalURL)
	id := shorturl
	urldb[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shorturl,
		CreationTime: time.Now(),
	}
	return shorturl
}

func geturl(id string) (URL, error) {
	url, ok := urldb[id]
	if !ok {
		return URL{}, errors.New("not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getmethod")
}

func shortenurl(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "eror", http.StatusBadRequest)
		return
	}

	shorturl := cretaeurl(data.URL)
	fmt.Println(w, shorturl)

	response := struct {
		ShortURL string `json:"shorturl"`
	}{ShortURL: shorturl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := geturl(id)
	if err != nil {
		http.Error(w, "invalid", http.StatusFound)
		return

	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	fmt.Println("starting the server")

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", shortenurl)
	http.HandleFunc("/redirect/", redirectURL)

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("errors", err)
	}

}
