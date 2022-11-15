package downloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	ErrDownloadServerNotFound      = errors.New("DownloadError: server not found")
	ErrDownloadIncorrectStatusCode = errors.New("DownloadError: incorrect status code")
	ErrDownload                    = errors.New("DownloadError")
)

type Downloader struct {
	c *http.Client
}

func NewDownloader(c *http.Client) Downloader {
	return Downloader{
		c: c,
	}
}

func (d Downloader) DownloadImg(req http.Request, file *os.File) error {
	resp, err := d.c.Do(&req)
	if err != nil {
		log.Println(err)
		return ErrDownloadServerNotFound
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	if resp.StatusCode >= 400 {
		log.Println("Server response code: ", resp.StatusCode)
		return ErrDownloadIncorrectStatusCode
	}

	if err := file.Truncate(0); err != nil {
		log.Println(err)
		return ErrDownload
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Seek err: ", err)
		return ErrDownload
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return ErrDownload
	}

	return nil
}
