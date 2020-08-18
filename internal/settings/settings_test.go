package settings

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	environment = struct {
		port      string // ~IMAGE_PREVIEWER_PORT
		cacheSize string // ~IMAGE_PREVIEWER_CACHE_SIZE
		minWidth  string // ~IMAGE_PREVIEWER_MIN_WIDTH
		minHeight string // ~IMAGE_PREVIEWER_MIN_HEIGHT
		maxWidth  string // ~IMAGE_PREVIEWER_MAX_WIDTH
		maxHeight string // ~IMAGE_PREVIEWER_MAX_HEIGHT
	}

	testCase = struct {
		name     string
		env      environment
		expected *Settings
		err      error
	}
)

var testCases = []testCase{
	{
		name:     "positive",
		env:      environment{"8080", "5", "50", "50", "2000", "2000"},
		expected: &Settings{8080, 5, 50, 50, 2000, 2000},
		err:      nil,
	},
	{
		name:     "canNotParsePort",
		env:      environment{"port", "5", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_PORT")),
	},
	{
		name:     "canNotParseCacheSize",
		env:      environment{"8080", "cacheSize", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_CACHE_SIZE")),
	},
	{
		name:     "canNotParseMinWidth",
		env:      environment{"8080", "5", "minWidth", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_MIN_WIDTH")),
	},
	{
		name:     "canNotParseMinHeight",
		env:      environment{"8080", "5", "50", "minHeight", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_MIN_HEIGHT")),
	},
	{
		name:     "canNotParseMaxWidth",
		env:      environment{"8080", "5", "50", "50", "maxWidth", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_MAX_WIDTH")),
	},
	{
		name:     "canNotParseMaxHeight",
		env:      environment{"8080", "5", "50", "50", "2000", "maxHeight"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			errors.New("can not parse IMAGE_PREVIEWER_MAX_HEIGHT")),
	},
	{
		name:     "portBoundsLeft",
		env:      environment{"-1", "5", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_PORT value must be in range [%d, %d]", minPort, maxPort)),
	},
	{
		name:     "portBoundsRight",
		env:      environment{"65536", "5", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_PORT value must be in range [%d, %d]", minPort, maxPort)),
	},
	{
		name:     "cacheSizeBoundsLeft",
		env:      environment{"8080", "0", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_CACHE_SIZE value must be in range [%d, %d]", minCacheSize, maxCacheSize)),
	},
	{
		name:     "cacheSizeBoundsRight",
		env:      environment{"8080", "10001", "50", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_CACHE_SIZE value must be in range [%d, %d]", minCacheSize, maxCacheSize)),
	},
	{
		name:     "minWidthBoundsLeft",
		env:      environment{"8080", "5", "0", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MIN_WIDTH value must be in range [%d, %d]", minMinWidth, maxMinWidth)),
	},
	{
		name:     "minWidthBoundsRight",
		env:      environment{"8080", "5", "1001", "50", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MIN_WIDTH value must be in range [%d, %d]", minMinWidth, maxMinWidth)),
	},
	{
		name:     "minHeightBoundsLeft",
		env:      environment{"8080", "5", "50", "0", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MIN_HEIGHT value must be in range [%d, %d]", minMinHeight, maxMinHeight)),
	},
	{
		name:     "minHeightBoundsRight",
		env:      environment{"8080", "5", "50", "1001", "2000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MIN_HEIGHT value must be in range [%d, %d]", minMinHeight, maxMinHeight)),
	},
	{
		name:     "maxWidthBoundsLeft",
		env:      environment{"8080", "5", "50", "50", "1000", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MAX_WIDTH value must be in range [%d, %d]", minMaxWidth, maxMaxWidth)),
	},
	{
		name:     "maxWidthBoundsRight",
		env:      environment{"8080", "5", "50", "50", "10001", "2000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MAX_WIDTH value must be in range [%d, %d]", minMaxWidth, maxMaxWidth)),
	},
	{
		name:     "maxHeightBoundsLeft",
		env:      environment{"8080", "5", "50", "50", "2000", "1000"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MAX_HEIGHT value must be in range [%d, %d]", minMaxHeight, maxMaxHeight)),
	},
	{
		name:     "maxHeightBoundsRight",
		env:      environment{"8080", "5", "50", "50", "2000", "10001"},
		expected: &Settings{0, 0, 0, 0, 0, 0},
		err: fmt.Errorf("%s: %w", ErrCanNotGetSettings,
			fmt.Errorf("IMAGE_PREVIEWER_MAX_HEIGHT value must be in range [%d, %d]", minMaxHeight, maxMaxHeight)),
	},
}

func TestParseEnv(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := new(Settings)
			setEnv(tc.env) //nolint:go-lint // using of "tc" in anonimous function
			defer unsetEnv()
			err := settings.ParseEnv()
			require.Equal(t, tc.err, err)           //nolint:go-lint
			require.Equal(t, tc.expected, settings) //nolint:go-lint
		})
	}
}

func setEnv(values environment) {
	os.Setenv("IMAGE_PREVIEWER_PORT", values.port)
	os.Setenv("IMAGE_PREVIEWER_CACHE_SIZE", values.cacheSize)
	os.Setenv("IMAGE_PREVIEWER_MIN_WIDTH", values.minWidth)
	os.Setenv("IMAGE_PREVIEWER_MIN_HEIGHT", values.minHeight)
	os.Setenv("IMAGE_PREVIEWER_MAX_WIDTH", values.maxWidth)
	os.Setenv("IMAGE_PREVIEWER_MAX_HEIGHT", values.maxHeight)
}

func unsetEnv() {
	os.Unsetenv("IMAGE_PREVIEWER_PORT")
	os.Unsetenv("IMAGE_PREVIEWER_CACHE_SIZE")
	os.Unsetenv("IMAGE_PREVIEWER_MIN_WIDTH")
	os.Unsetenv("IMAGE_PREVIEWER_MIN_HEIGHT")
	os.Unsetenv("IMAGE_PREVIEWER_MAX_WIDTH")
	os.Unsetenv("IMAGE_PREVIEWER_MAX_HEIGHT")
}
