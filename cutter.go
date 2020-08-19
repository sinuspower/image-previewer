package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageCutter interface {
	LoadImage() ([]byte, error)
	Cut([]byte) ([]byte, error)
}

type Cutter struct {
	width  int
	height int
	url    string
	client *http.Client
}

var ErrCanNotParsePath = errors.New("can not parse path")

func NewCutter(path string) (ImageCutter, error) {
	width, height, url, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	return &Cutter{
		width:  width,
		height: height,
		url:    url,
		client: &http.Client{},
	}, nil
}

func (c *Cutter) LoadImage() ([]byte, error) {
	// load source image from cache
	image, ok, err := cache.GetFile(c.url)
	if err != nil {
		log.Println("[WARN] can not get source image from cache:", err)
	}
	if ok {
		log.Println("[INFO] get source image from cache")

		return image, nil
	}
	// --------------------------

	rq, err := http.NewRequestWithContext(context.Background(), "GET", c.url, nil)
	if err != nil {
		return nil, err
	}

	log.Println("[INFO] send request to", c.url)
	rs, err := c.client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	log.Println("[INFO] get response from", c.url)
	bytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	err = cache.PutFile(c.url, bytes)
	if err != nil {
		log.Println("[WARN] can put source image into cache:", err)
	} else {
		log.Println("[INFO] put source image into cache")
	}

	return bytes, nil
}

func (c *Cutter) Cut(source []byte) ([]byte, error) {
	image, _, err := image.Decode(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	preview := imaging.Fill(image, c.width, c.height, imaging.Center, imaging.Lanczos)
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, preview, nil)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// parsePath returns width, height and URL from input string like /fill/300/200/{URL}.
func parsePath(path string) (int, int, string, error) {
	parts := strings.SplitN(path, "/", 5)

	if len(parts) < 5 {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath,
			errors.New("missing expected elements in URL"))
	}

	width, err := getWidth(parts[2])
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath, err)
	}

	height, err := getHeight(parts[3])
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath, err)
	}

	url, err := getURL(parts[4])
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath, err)
	}

	return width, height, url, nil
}

func getWidth(source string) (int, error) {
	width, err := strconv.Atoi(source)
	if err != nil {
		return 0, errors.New("can not get width")
	}

	min := settings.GetMinWidth()
	max := settings.GetMaxWidth()
	if width < min || width > max {
		return 0, fmt.Errorf("width value must be in range [%d, %d]", min, max)
	}

	return width, nil
}

func getHeight(source string) (int, error) {
	height, err := strconv.Atoi(source)
	if err != nil {
		return 0, errors.New("can not get height")
	}

	min := settings.GetMinHeight()
	max := settings.GetMaxHeight()
	if height < min || height > max {
		return 0, fmt.Errorf("height value must be in range [%d, %d]", min, max)
	}

	return height, nil
}

func getURL(source string) (string, error) {
	if !strings.HasSuffix(source, "jpg") && !strings.HasSuffix(source, "jpeg") {
		return "", errors.New("file extension must be jpg or jpeg")
	}

	if !strings.HasPrefix(source, "http://") {
		source = "http://" + source
	}

	return source, nil
}
