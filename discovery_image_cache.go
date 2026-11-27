package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	discoveryImageMaxFiles = 640
	discoveryImageMaxBytes = int64(32 << 20)
	discoveryImageMaxAge   = 60 * 24 * time.Hour
	discoveryImageMaxSize  = 2 << 20
	discoveryImageTouchAge = 12 * time.Hour
)

// Discovery uses larger images than the library thumbnails and has its own budget.
type discoveryImageCache struct {
	mu          sync.Mutex
	files       int
	bytes       int64
	lastCleanup time.Time
}

func discoveryImageDir(libraryPath string) string {
	return filepath.Join(appCacheRoot(filepath.Dir(libraryPath)), "imagenes_descubrir")
}

func discoveryImageType(data []byte) string {
	switch mime := http.DetectContentType(data); mime {
	case "image/jpeg", "image/png", "image/webp":
		return mime
	default:
		return ""
	}
}

func posterDiskUsage(dir string) (files int, bytes int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !validPosterPath("/"+entry.Name()) {
			continue
		}
		if info, err := entry.Info(); err == nil {
			files++
			bytes += info.Size()
		}
	}
	return files, bytes
}

func (c *discoveryImageCache) cleanup(dir string, maxFiles int, maxBytes int64, maxAge time.Duration) (int, error) {
	removed, err := prunePosterCache(dir, maxFiles, maxBytes, maxAge)
	if err == nil {
		c.files, c.bytes = posterDiskUsage(dir)
		c.lastCleanup = time.Now()
	}
	return removed, err
}

func (c *discoveryImageCache) maintain(dir string, clear bool) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if clear {
		return c.cleanup(dir, 0, 0, 0)
	}
	return c.cleanup(dir, discoveryImageMaxFiles, discoveryImageMaxBytes, discoveryImageMaxAge)
}

func (c *discoveryImageCache) get(dir, name string) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Since(c.lastCleanup) >= 24*time.Hour {
		_, _ = c.cleanup(dir, discoveryImageMaxFiles, discoveryImageMaxBytes, discoveryImageMaxAge)
	}
	path := filepath.Join(dir, name)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > discoveryImageMaxSize ||
		time.Since(info.ModTime()) > discoveryImageMaxAge {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	data, err := io.ReadAll(io.LimitReader(file, discoveryImageMaxSize+1))
	_ = file.Close()
	if err != nil || len(data) > discoveryImageMaxSize || discoveryImageType(data) == "" {
		return nil
	}
	if time.Since(info.ModTime()) >= discoveryImageTouchAge {
		now := time.Now()
		_ = os.Chtimes(path, now, now)
	}
	return data
}

func (c *discoveryImageCache) put(dir, name string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lastCleanup.IsZero() {
		if _, err := c.cleanup(dir, discoveryImageMaxFiles, discoveryImageMaxBytes, discoveryImageMaxAge); err != nil {
			return
		}
	}
	path := filepath.Join(dir, name)
	existing := func() (int, int64) {
		if info, err := os.Lstat(path); err == nil && info.Mode().IsRegular() {
			return 1, info.Size()
		}
		return 0, 0
	}
	oldFiles, oldBytes := existing()
	size := int64(len(data))
	if c.files-oldFiles+1 > discoveryImageMaxFiles || c.bytes-oldBytes+size > discoveryImageMaxBytes {
		// Evict in batches so a full cache does not rescan the directory for every miss.
		if _, err := c.cleanup(dir, discoveryImageMaxFiles*7/8, discoveryImageMaxBytes*7/8, discoveryImageMaxAge); err != nil {
			return
		}
		oldFiles, oldBytes = existing()
	}
	if c.files-oldFiles+1 > discoveryImageMaxFiles || c.bytes-oldBytes+size > discoveryImageMaxBytes {
		return
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return
	}
	if err := replaceLibraryFile(path, data); err == nil {
		c.files += 1 - oldFiles
		c.bytes += size - oldBytes
	}
}
