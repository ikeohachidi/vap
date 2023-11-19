package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"
	"github.com/valyala/fasthttp"
)

const GH_API_BASE = "https://api.github.com"

type FileMetaData struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

type Repo struct {
	username string
	title    string
}

type PageContent struct {
	repository Repo
}

func (h *PageContent) fetchPage(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	repoURL := fmt.Sprintf("%v/repos/%v/%v/contents/%v", GH_API_BASE, h.repository.username, h.repository.title, path)

	resp, err := http.Get(repoURL)
	if err != nil {
		log.Fatal(err)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var fileMeta FileMetaData
	err = json.Unmarshal(bodyByte, &fileMeta)
	if err != nil {
		log.Fatal(err)
	}

	body, err := base64.StdEncoding.DecodeString(fileMeta.Content)
	if err != nil {
		log.Fatal(err)
	}

	extension := filepath.Ext(fileMeta.DownloadURL)

	mediaType := mime.TypeByExtension(extension)
	ctx.SetContentType(mediaType)
	ctx.Write(body)
}

func main() {
	repo := os.Args[1]

	rURL, err := url.Parse(repo)
	if err != nil {
		log.Fatal(err)
	}

	// path should come in the form <username>/<repo_title>
	splitStr := strings.Split(rURL.Path, "/")
	if len(splitStr) < 2 {
		log.Fatal("Not enough info in url")
	}

	page := &PageContent{
		repository: Repo{
			username: splitStr[1],
			title:    splitStr[2],
		},
	}

	log.Println("Opening on port 8080")
	browser.OpenURL("http://localhost:8080/index.html")
	err = fasthttp.ListenAndServe(":8080", page.fetchPage)
	if err != nil {
		log.Fatal(err)
	}
}
