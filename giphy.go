package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var giphySearchURL = "http://api.giphy.com/v1/gifs/search"

type GiphyImageResponse struct {
	EmbedURL string `json:"embed_url"`
	URL      string `json:"url"`
}

type giphyImageHandler struct {
	APIKey     string
	AprilFools bool
}

func (g *giphyImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var embedURL string
	var imageURL string

	keys, ok := r.URL.Query()["keyword"]
	if !ok || len(keys) < 1 {
		log.Println("URL parameter 'keyword' is missing")
		http.Error(w, "Request doesn't have 'keyword' as a URL parameter", http.StatusNotFound)
		return
	}

	// If April fools feature flag is set then return
	// static images
	if g.AprilFools {
		json.NewEncoder(w).Encode(GiphyImageResponse{
			EmbedURL: "https://i.ytimg.com/vi/JrQkgLLL9XQ/hqdefault.jpg",
			URL:      "https://i.ytimg.com/vi/JrQkgLLL9XQ/hqdefault.jpg",
		})
		return
	}

	url := fmt.Sprintf("%s?q=%s&api_key=%s&limit=5", giphySearchURL, keys[0], g.APIKey)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make a request to giphy: %v", err)
		http.Error(w, "Failed to make a request to giphy", http.StatusInternalServerError)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		http.Error(w, "Failed to read response from giphy", http.StatusInternalServerError)
	}
	if res.StatusCode != 200 {
		log.Printf("error response '%d' -> %s", res.StatusCode, string(body))
	}
	var s SearchResponse
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("Failed to decode json: %v", err)
		http.Error(w, "Failed to decode json response", http.StatusInternalServerError)
	}

	if len(s.Data) > 0 {
		embedURL = s.Data[0].EmbedURL
		imageURL = s.Data[0].URL
	}

	response := GiphyImageResponse{
		EmbedURL: embedURL,
		URL:      imageURL,
	}
	json.NewEncoder(w).Encode(response)
	return
}

func GiphyImageHandler(apiKey string, aprilFools bool) http.Handler {
	return &giphyImageHandler{
		APIKey:     apiKey,
		AprilFools: aprilFools,
	}
}
