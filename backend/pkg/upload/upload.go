// Package upload provides small helpers for validating and storing user-uploaded
// image files on local disk (shared volume served by nginx as /uploads).
package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Allowed image MIME types for person photos.
var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// DetectImageMIME sniffs the content type from the first bytes of data.
// Returns ErrUnsupported when the type is not an allowed image.
func DetectImageMIME(header []byte) (mimeType, ext string, err error) {
	mimeType = http.DetectContentType(header)
	// jpeg sniff sometimes returns image/jpeg; strip parameters if any.
	mimeType = strings.Split(mimeType, ";")[0]
	ext, ok := allowedMIME[mimeType]
	if !ok {
		return "", "", fmt.Errorf("unsupported image type %q", mimeType)
	}
	return mimeType, ext, nil
}

// Save writes data to dir/relativePath (creating parent dirs). relativePath
// must be a clean path with no ".." segments.
func Save(dir, relativePath string, data []byte) error {
	rel := filepath.Clean(relativePath)
	if rel == "." || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("invalid relative path")
	}
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	tmp := full + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write temp: %w", err)
	}
	if err := os.Rename(tmp, full); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

// Remove deletes dir/relativePath if it exists. Empty relativePath is a no-op.
func Remove(dir, relativePath string) error {
	if strings.TrimSpace(relativePath) == "" {
		return nil
	}
	rel := filepath.Clean(relativePath)
	if rel == "." || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("invalid relative path")
	}
	full := filepath.Join(dir, rel)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ReadLimited reads at most maxBytes+1 from r so callers can detect overflow.
func ReadLimited(r io.Reader, maxBytes int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file too large")
	}
	return data, nil
}
