package integration_test

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
}

func (s *IntegrationSuite) TestImageInCache() {
	res1, err := s.sendRequest(200, 500, "http://web/testdata/in_image.jpg")
	defer func() {
		if err := res1.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	s.Require().Equal(res1.StatusCode, http.StatusOK)

	res, err := s.sendRequest(200, 500, "http://web/testdata/in_image.jpg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	s.Require().Equal(res.StatusCode, http.StatusOK)

	val := res.Header.Get("Image-From-Cache")

	s.Require().Equal(val, "1")

	sumRespFile, err := s.getMD5SumString(res.Body)

	s.Require().NoError(err)

	file, err := os.Open("./../../testdata/out_cache_image.jpg")

	s.Require().NoError(err)

	sumTestFile, err := s.getMD5SumString(file)

	s.Require().NoError(err)

	s.Require().Equal(sumRespFile, sumTestFile)
}

func (s *IntegrationSuite) TestImageOnServer() {
	res, err := s.sendRequest(300, 400, "http://web/testdata/in_image.jpg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	s.Require().Equal(res.StatusCode, http.StatusOK)

	sumRespFile, err := s.getMD5SumString(res.Body)

	s.Require().NoError(err)

	file, err := os.Open("./../../testdata/out_image.jpg")

	s.Require().NoError(err)

	sumTestFile, err := s.getMD5SumString(file)

	s.Require().NoError(err)

	s.Require().Equal(sumRespFile, sumTestFile)
}

func (s *IntegrationSuite) TestServerNotFound() {
	res, err := s.sendRequest(200, 500, "http://web:81/testdata/in_image.jpg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	body, err := io.ReadAll(res.Body)

	s.Require().NoError(err)

	s.Require().Equal(string(body), "DownloadError: server not found\n")
}

func (s *IntegrationSuite) TestImageNotFoundOnServer() {
	res, err := s.sendRequest(200, 500, "http://web/testdata/image.jpg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	body, err := io.ReadAll(res.Body)

	s.Require().NoError(err)

	s.Require().Equal(string(body), "DownloadError: incorrect status code\n")
}

func (s *IntegrationSuite) TestFileNotImage() {
	res, err := s.sendRequest(200, 500, "http://web/testdata/in_run.exe")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	body, err := io.ReadAll(res.Body)

	s.Require().NoError(err)

	s.Require().Equal(string(body), "url incorrect\n")
}

func (s *IntegrationSuite) TestIncorrectSizeImageOnServer() {
	res, err := s.sendRequest(200, 500, "http://web/testdata/in_small_image.jpeg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	body, err := io.ReadAll(res.Body)

	s.Require().NoError(err)

	s.Require().Equal(string(body), "incorrect image size\n")
}

func (s *IntegrationSuite) TestIncorrectSizeImageForResize() {
	res, err := s.sendRequest(1, 100000, "http://web/testdata/in_small_image.jpeg")
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	s.Require().NoError(err)

	body, err := io.ReadAll(res.Body)

	s.Require().NoError(err)

	s.Require().Equal(string(body), "incorrect size\n")
}

func (s *IntegrationSuite) sendRequest(width, height int, urlForImage string) (*http.Response, error) {
	url := fmt.Sprintf("http://resizer/fill/%d/%d/%s", width, height, urlForImage)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (s *IntegrationSuite) getMD5SumString(f io.Reader) (string, error) {
	file1Sum := md5.New()
	if _, err := io.Copy(file1Sum, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
