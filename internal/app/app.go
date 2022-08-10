package app

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/antonbaks/otusProjectJob/internal/lru"
)

type App struct {
	cache lru.Cache
	s     Storage
	r     Resizer
	d     Downloader
}

type Resizer interface {
	ResizeImage(w, h uint, file *os.File) error
}

type Downloader interface {
	DownloadImg(req http.Request, file *os.File) error
}

type Storage interface {
	Open(url, width, height string) (*os.File, error)
	Create(url, width, height string) (*os.File, error)
	FileName(url, width, height string) string
	CreateUploadDir() error
}

func NewApp(c lru.Cache, s Storage, d Downloader, r Resizer) App {
	return App{
		cache: c,
		s:     s,
		d:     d,
		r:     r,
	}
}

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting dispatch url: ", r.URL.Path)
	dto, err := urlDispatcher(r.URL.Path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err = h.s.CreateUploadDir(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("Checking in cache")
	imageFromCache := false
	cacheKey := lru.Key(h.s.FileName(dto.imgURL.String(), strconv.Itoa(dto.width), strconv.Itoa(dto.height)))
	if _, ok := h.cache.Get(cacheKey); ok {
		log.Println("Find in cache")
		file, err := h.s.Open(dto.imgURL.String(), strconv.Itoa(dto.width), strconv.Itoa(dto.height))
		defer closeFile(file)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		imageFromCache = true

		if err := imgResponse(w, file, imageFromCache); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
		}

		return
	}

	log.Println("Create request for download")
	req, err := http.NewRequestWithContext(context.Background(), "GET", dto.imgURL.String(), nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	req.Header = r.Header

	log.Println("Create empty image")
	file, err := h.s.Create(dto.imgURL.String(), strconv.Itoa(dto.width), strconv.Itoa(dto.height))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer closeFile(file)

	log.Println("Downloading image")
	if err = h.d.DownloadImg(*req, file); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("Resizing image")
	if err := h.r.ResizeImage(uint(dto.width), uint(dto.height), file); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("Add image in cache")
	h.cache.Set(cacheKey, "")

	if err := imgResponse(w, file, imageFromCache); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Println(err)
		return
	}
}

func imgResponse(w http.ResponseWriter, file *os.File, imageFromCache bool) error {
	log.Println("Create image response")

	w.Header().Set("Content-Type", "image/jpeg")
	if imageFromCache {
		w.Header().Set("Image-From-Cache", "1")
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}
