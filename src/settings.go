package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const ( // foolproof
	minPort      = 0
	maxPort      = 65535
	minCacheSize = 1
	maxCacheSize = 10000
	minMinWidth  = 1
	maxMinWidth  = 1000
	minMinHeight = 1
	maxMinHeight = 1000
	minMaxWidth  = maxMinWidth + 1
	maxMaxWidth  = 10000
	minMaxHeight = maxMinHeight + 1
	maxMaxHeight = 10000
)

type Settings struct {
	port      int // ~IMAGE_PREVIEWER_PORT
	cacheSize int // ~IMAGE_PREVIEWER_CACHE_SIZE
	minWidth  int // ~IMAGE_PREVIEWER_MIN_WIDTH
	minHeight int // ~IMAGE_PREVIEWER_MIN_HEIGHT
	maxWidth  int // ~IMAGE_PREVIEWER_MAX_WIDTH
	maxHeight int // ~IMAGE_PREVIEWER_MAX_HEIGHT
}

var ErrCanNotGetSettings = errors.New("can not get settings")

func (s *Settings) ParseEnv() error {
	port, err := parseIntVar("IMAGE_PREVIEWER_PORT", minPort, maxPort)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.port = port

	cacheSize, err := parseIntVar("IMAGE_PREVIEWER_CACHE_SIZE", minCacheSize, maxCacheSize)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.cacheSize = cacheSize

	minWidth, err := parseIntVar("IMAGE_PREVIEWER_MIN_WIDTH", minMinWidth, maxMinWidth)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.minWidth = minWidth

	minHeight, err := parseIntVar("IMAGE_PREVIEWER_MIN_HEIGHT", minMinHeight, maxMinHeight)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.minHeight = minHeight

	maxWidth, err := parseIntVar("IMAGE_PREVIEWER_MAX_WIDTH", minMaxWidth, maxMaxWidth)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.maxWidth = maxWidth

	maxHeight, err := parseIntVar("IMAGE_PREVIEWER_MAX_HEIGHT", minMaxHeight, maxMaxHeight)
	if err != nil {
		s.Reset()

		return fmt.Errorf("%s: %w", ErrCanNotGetSettings, err)
	}
	s.maxHeight = maxHeight

	return nil
}

func (s *Settings) GetPort() int {
	return s.port
}

func (s *Settings) GetCacheSize() int {
	return s.cacheSize
}

func (s *Settings) GetMinWidth() int {
	return s.minWidth
}

func (s *Settings) GetMinHeight() int {
	return s.minHeight
}

func (s *Settings) GetMaxWidth() int {
	return s.maxWidth
}

func (s *Settings) GetMaxHeight() int {
	return s.maxHeight
}

func (s *Settings) Reset() {
	s.port, s.cacheSize, s.minWidth, s.minHeight, s.maxWidth, s.maxHeight = 0, 0, 0, 0, 0, 0
}

func parseIntVar(name string, min int, max int) (int, error) {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return 0, fmt.Errorf("can not parse %s", name)
	}

	if value < min || value > max {
		return 0, fmt.Errorf("%s value must be in range [%d, %d]", name, min, max)
	}

	return value, nil
}
