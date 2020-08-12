package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePath(t *testing.T) { //nolint:go-lint // (function too long)
	type expected = struct {
		width  uint16
		height uint16
		url    string
		err    string
	}
	type testCase = struct {
		in       string
		expected expected
	}

	positive := []testCase{
		{
			in:       "/fill/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
			expected: expected{300, 200, "http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg", ""},
		},
		{
			in:       "/fill/300/200/http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
			expected: expected{300, 200, "http://www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg", ""},
		},
		{
			in:       "/fill/100/100/path/path/image.jpg",
			expected: expected{100, 100, "http://path/path/image.jpg", ""},
		},
		{
			in:       "/fill/100/100/path/path/image.jpeg",
			expected: expected{100, 100, "http://path/path/image.jpeg", ""},
		},
	}

	negative := []testCase{
		{
			in:       "/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg",
			expected: expected{0, 0, "", "can not parse path: strconv.ParseUint: parsing \"www.audubon.org\": invalid syntax"},
		},
		{
			in:       "bad",
			expected: expected{0, 0, "", "can not parse path: missing expected elements"},
		},
		{
			in:       "",
			expected: expected{0, 0, "", "can not parse path: missing expected elements"},
		},
		{
			in:       "/fill/100/100/path/path/song.flac",
			expected: expected{0, 0, "", "file extension must be jpg or jpeg"},
		},
		{
			in:       "/fill/100/100/path/path/doc.pdf",
			expected: expected{0, 0, "", "file extension must be jpg or jpeg"},
		},
	}

	for _, tc := range positive {
		t.Run(tc.in, func(t *testing.T) {
			width, height, url, err := parsePath(tc.in)  //nolint:go-lint //using of "tc" in anonimous function
			require.Equal(t, tc.expected.width, width)   //nolint:go-lint
			require.Equal(t, tc.expected.height, height) //nolint:go-lint
			require.Equal(t, tc.expected.url, url)       //nolint:go-lint
			require.NoError(t, err)
		})
	}

	for _, tc := range negative {
		t.Run(tc.in, func(t *testing.T) {
			width, height, url, err := parsePath(tc.in)  //nolint:go-lint
			require.Equal(t, tc.expected.width, width)   //nolint:go-lint
			require.Equal(t, tc.expected.height, height) //nolint:go-lint
			require.Equal(t, tc.expected.url, url)       //nolint:go-lint
			require.EqualError(t, err, tc.expected.err)  //nolint:go-lint
		})
	}
}
