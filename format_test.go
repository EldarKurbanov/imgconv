package imgconv

import (
	"bytes"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestSetFormat(t *testing.T) {
	if _, err := setFormat("Jpg"); err != nil {
		t.Error("Failed to set format")
	}
	if _, err := setFormat("txt"); err == nil {
		t.Error("set txt format want error")
	}
}

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/video-001.png")
	if err != nil {
		t.Error(err)
		return
	}
	if _, err := Decode(f); err != nil {
		t.Error("Failed to decode")
	}
	if _, err := Decode(bytes.NewBufferString("Hello")); err == nil {
		t.Error("Decode string want error")
	}
}

func TestEncode(t *testing.T) {
	testCase := []FormatOption{
		{Format: JPEG, EncodeOption: []EncodeOption{JPEGQuality(75)}},
		{Format: PNG, EncodeOption: []EncodeOption{PNGCompressionLevel(png.DefaultCompression)}},
		{Format: GIF, EncodeOption: []EncodeOption{GIFNumColors(256), GIFDrawer(draw.FloydSteinberg), GIFQuantizer(nil)}},
		{Format: TIFF},
		{Format: BMP},
	}

	// Read the image.
	m0, err := Open("testdata/video-001.png")
	if err != nil {
		t.Error(err)
		return
	}
	for _, tc := range testCase {
		// Encode the image.
		var buf bytes.Buffer
		fo, err := setFormat(formatExts[tc.Format], tc.EncodeOption...)
		if err != nil {
			t.Error(tc, err)
			continue
		}
		if err := Write(m0, &buf, fo); err != nil {
			t.Error(formatExts[fo.Format], err)
			continue
		}
		// Decode the image.
		b, _ := ioutil.ReadAll(&buf)
		m1, err := decode(bytes.NewBuffer(b), fo.Format)
		if err != nil {
			t.Error(formatExts[fo.Format], err)
			continue
		}
		if m0.Bounds() != m1.Bounds() {
			t.Errorf("bounds differ: %v and %v", m0.Bounds(), m1.Bounds())
			continue
		}
	}
}

func TestOpenSave(t *testing.T) {
	if _, err := Open("/dev/null"); err == nil {
		t.Error("Open invalid image want error")
	}
	img, err := Open("testdata/video-001.png")
	if err != nil {
		t.Error("Fail to open image", err)
		return
	}
	if err := Save(img, "testdata/video-001.jpg", defaultFormat); err != nil {
		t.Error("Fail to save image", err)
		return
	}
	if err := os.Remove("testdata/video-001.jpg"); err != nil {
		t.Error(err)
	}
}
