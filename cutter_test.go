package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	internal_settings "github.com/sinuspower/image-previewer/internal/settings"
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

	expected = struct {
		width  int
		height int
		url    string
		err    string
	}

	parsePathTestCase = struct {
		name     string
		in       string
		expected expected
	}
)

func TestParsePath(t *testing.T) { //nolint:go-lint // function is too long
	env := environment{"8080", "5", "50", "50", "2000", "2000"}
	setEnv(env)
	settings = new(internal_settings.Settings)
	err := settings.ParseEnv()
	require.NoError(t, err)

	parsePathPositive, parsePathNegative := getParsePathTestCases()

	for _, tc := range parsePathPositive {
		t.Run(tc.name, func(t *testing.T) {
			width, height, url, err := parsePath(tc.in)  //nolint:go-lint // using of "tc" in anonimous function
			require.Equal(t, tc.expected.width, width)   //nolint:go-lint
			require.Equal(t, tc.expected.height, height) //nolint:go-lint
			require.Equal(t, tc.expected.url, url)       //nolint:go-lint
			require.NoError(t, err)
		})
	}

	for _, tc := range parsePathNegative {
		t.Run(tc.name, func(t *testing.T) {
			width, height, url, err := parsePath(tc.in)  //nolint:go-lint
			require.Equal(t, tc.expected.width, width)   //nolint:go-lint
			require.Equal(t, tc.expected.height, height) //nolint:go-lint
			require.Equal(t, tc.expected.url, url)       //nolint:go-lint
			require.EqualError(t, err, tc.expected.err)  //nolint:go-lint
		})
	}
}

func TestCut(t *testing.T) {
	source1024x504, err := readFile("test/testdata/_gopher_original_1024x504.jpg")
	require.NoError(t, err)

	exp1024x504, err := readFile("test/testdata/gopher_1024x504.jpg")
	require.NoError(t, err)

	exp50x50, err := readFile("test/testdata/gopher_50x50.jpg")
	require.NoError(t, err)

	exp200x700, err := readFile("test/testdata/gopher_200x700.jpg")
	require.NoError(t, err)

	exp256x126, err := readFile("test/testdata/gopher_256x126.jpg")
	require.NoError(t, err)

	exp333x666, err := readFile("test/testdata/gopher_333x666.jpg")
	require.NoError(t, err)

	exp500x500, err := readFile("test/testdata/gopher_500x500.jpg")
	require.NoError(t, err)

	exp1024x252, err := readFile("test/testdata/gopher_1024x252.jpg")
	require.NoError(t, err)

	exp2000x1000, err := readFile("test/testdata/gopher_2000x1000.jpg")
	require.NoError(t, err)

	type testCase = struct {
		name     string
		path     string
		expected []byte
	}

	testCases := []testCase{
		{"1024x504", "/fill/1024/504/www.testcut.com/source.jpg", exp1024x504},
		{"50x50", "/fill/50/50/www.testcut.com/source.jpg", exp50x50},
		{"200x700", "/fill/200/700/www.testcut.com/source.jpg", exp200x700},
		{"256x126", "/fill/256/126/www.testcut.com/source.jpg", exp256x126},
		{"333x666", "/fill/333/666/www.testcut.com/source.jpg", exp333x666},
		{"500x500", "/fill/500/500/www.testcut.com/source.jpg", exp500x500},
		{"1024x252", "/fill/1024/252/www.testcut.com/source.jpg", exp1024x252},
		{"2000x1000", "/fill/2000/1000/www.testcut.com/source.jpg", exp2000x1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cutter, err := NewCutter(tc.path) //nolint:go-lint
			require.NoError(t, err)
			actual, err := cutter.Cut(source1024x504)
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual) //nolint:go-lint
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

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func getParsePathTestCases() ([]parsePathTestCase, []parsePathTestCase) { //nolint:go-lint // function is too long
	return []parsePathTestCase{
			{
				name:     "positiveWithoutHTTP",
				in:       "/fill/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
				expected: expected{300, 200, "http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg", ""},
			},
			{
				name:     "positiveWithHTTP",
				in:       "/fill/300/200/http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
				expected: expected{300, 200, "http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg", ""},
			},
			{
				name:     "positiveJPG",
				in:       "/fill/100/100/path/path/image.jpg",
				expected: expected{100, 100, "http://path/path/image.jpg", ""},
			},
			{
				name:     "positiveJPEG",
				in:       "/fill/100/100/path/path/image.jpeg",
				expected: expected{100, 100, "http://path/path/image.jpeg", ""},
			},
		}, []parsePathTestCase{
			{
				name:     "withoutFirstPathPart",
				in:       "/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
				expected: expected{0, 0, "", "can not parse path: can not get height"},
			},
			{
				name:     "oneWordPath",
				in:       "bad",
				expected: expected{0, 0, "", "can not parse path: missing expected elements in URL"},
			},
			{
				name:     "emptyPath",
				in:       "",
				expected: expected{0, 0, "", "can not parse path: missing expected elements in URL"},
			},
			{
				name:     "flacFile",
				in:       "/fill/100/100/path/path/song.flac",
				expected: expected{0, 0, "", "can not parse path: file extension must be jpg or jpeg"},
			},
			{
				name:     "pdfFile",
				in:       "/fill/100/100/path/path/doc.pdf",
				expected: expected{0, 0, "", "can not parse path: file extension must be jpg or jpeg"},
			},
			{
				name:     "canNotGetWidth",
				in:       "/fill/width/100/path/path/img.jpg",
				expected: expected{0, 0, "", "can not parse path: can not get width"},
			},
			{
				name:     "canNotGetHeight",
				in:       "/fill/100/height/path/path/img.jpg",
				expected: expected{0, 0, "", "can not parse path: can not get height"},
			},
			{
				name: "widthBoundsLeft",
				in:   fmt.Sprintf("/fill/%d/100/path/path/img.jpg", settings.GetMinWidth()-1),
				expected: expected{0, 0, "", fmt.Sprintf("can not parse path: width value must be in range [%d, %d]",
					settings.GetMinWidth(), settings.GetMaxWidth())},
			},
			{
				name: "widthBoundsRight",
				in:   fmt.Sprintf("/fill/%d/100/path/path/img.jpg", settings.GetMaxWidth()+1),
				expected: expected{0, 0, "", fmt.Sprintf("can not parse path: width value must be in range [%d, %d]",
					settings.GetMinWidth(), settings.GetMaxWidth())},
			},
			{
				name: "heightBoundsLeft",
				in:   fmt.Sprintf("/fill/100/%d/path/path/img.jpg", settings.GetMinHeight()-1),
				expected: expected{0, 0, "", fmt.Sprintf("can not parse path: height value must be in range [%d, %d]",
					settings.GetMinHeight(), settings.GetMaxHeight())},
			},
			{
				name: "heightBoundsRight",
				in:   fmt.Sprintf("/fill/100/%d/path/path/img.jpg", settings.GetMaxHeight()+1),
				expected: expected{0, 0, "", fmt.Sprintf("can not parse path: height value must be in range [%d, %d]",
					settings.GetMinHeight(), settings.GetMaxHeight())},
			},
		}
}
