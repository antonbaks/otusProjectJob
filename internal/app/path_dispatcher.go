package app

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

var ErrIncorrectURL = errors.New("url incorrect")

const pattern = `/fill/(\d*)/(\d*)/(.*\.(jpg|JPG|jpeg|JPEG))`

type ResizerRequestDTO struct {
	width  int
	height int
	imgURL *url.URL
}

func urlDispatcher(path string) (ResizerRequestDTO, error) {
	dto := ResizerRequestDTO{}

	re := regexp.MustCompile(pattern)

	matched := re.MatchString(path)

	if !matched {
		return dto, ErrIncorrectURL
	}

	res := re.FindAllStringSubmatch(path, -1)

	width, err := strconv.Atoi(res[0][1])
	if err != nil {
		return dto, ErrIncorrectURL
	}
	dto.width = width

	height, err := strconv.Atoi(res[0][2])
	if err != nil {
		return dto, ErrIncorrectURL
	}
	dto.height = height

	uri, err := url.ParseRequestURI(res[0][3])
	if err != nil {
		return dto, ErrIncorrectURL
	}

	dto.imgURL = uri

	return dto, nil
}
