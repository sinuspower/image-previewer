package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ImageCutter interface {
	LoadImage() ([]byte, error)
	Cut([]byte) ([]byte, error)
}

type Cutter struct {
	width  uint16
	height uint16
	url    string
	client *http.Client
}

var (
	ErrCanNotParsePath    = errors.New("can not parse path")
	ErrWrongFileExtension = errors.New("file extension must be jpg or jpeg")
)

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
	rq, err := http.NewRequestWithContext(context.Background(), "GET", c.url, nil)
	if err != nil {
		return nil, err
	}

	rs, err := c.client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	bytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c *Cutter) Cut(source []byte) ([]byte, error) {
	return source, nil
}

// parsePath returns width, height and URL from input string like /fill/300/200/{URL}.
func parsePath(path string) (uint16, uint16, string, error) {
	parts := strings.SplitN(path, "/", 5)

	if len(parts) < 5 {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath,
			errors.New("missing expected elements"))
	}

	width, err := strconv.ParseUint(parts[2], 10, 16)
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath, err)
	}

	height, err := strconv.ParseUint(parts[3], 10, 16)
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %w", ErrCanNotParsePath, err)
	}

	url := parts[4]
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	if !strings.HasSuffix(url, "jpg") && !strings.HasSuffix(url, "jpeg") {
		return 0, 0, "", ErrWrongFileExtension
	}

	return uint16(width), uint16(height), url, nil
}
