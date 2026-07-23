package upload_test

import (
	"strings"
	"testing"

	"github.com/twkakbar/family-tree/backend/pkg/upload"
)

func TestDetectImageMIME_JPEG(t *testing.T) {
	// Minimal JPEG SOI + JFIF-ish bytes for DetectContentType.
	header := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46}
	mime, ext, err := upload.DetectImageMIME(header)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if mime != "image/jpeg" || ext != ".jpg" {
		t.Fatalf("got %s %s", mime, ext)
	}
}

func TestDetectImageMIME_RejectsText(t *testing.T) {
	_, _, err := upload.DetectImageMIME([]byte("hello world not an image"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unexpected err: %v", err)
	}
}
