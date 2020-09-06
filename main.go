package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("")
		os.Exit(1)
	}

	name := os.Args[1]
	err := downloadPoster(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Download poster for movie %s\n", name)
}

func searchMovie(title string) (map[string]interface{}, error) {
	apiKey := os.Getenv("OMDB_APIKEY")
	url := fmt.Sprintf("http://www.omdbapi.com?apikey=%s&t=%s", apiKey, title)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve movie info: %s", resp.Status)
	}

	data := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&data)

	return data, nil
}

func downloadPoster(name string) error {
	movie, err := searchMovie(name)
	if err != nil {
		return err
	}

	resp, err := http.Get(movie["Poster"].(string))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filepath.Base(movie["Poster"].(string)))
	err = ioutil.WriteFile(fmt.Sprintf("posters/%s%s", normalizeMovieName(name), ext), blob, 0755)
	if err != nil {
		return err
	}

	return nil
}

func normalizeMovieName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}
