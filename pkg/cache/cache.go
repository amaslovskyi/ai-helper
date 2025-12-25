package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/ai-helper/pkg/llm"
)

// Entry represents a cache entry
type Entry struct {
	Command   string    `json:"cmd"`
	Error     string    `json:"error"`
	Fix       string    `json:"fix"`
	Timestamp int64     `json:"timestamp"`
	Hits      int       `json:"hits"`
	LastUsed  int64     `json:"last_used"`
}

// Cache manages the response cache
type Cache struct {
	file    string
	entries map[string]*Entry
}

// NewCache creates a new cache
func NewCache(cacheFile string) (*Cache, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(cacheFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	c := &Cache{
		file:    cacheFile,
		entries: make(map[string]*Entry),
	}

	// Load existing cache
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	return c, nil
}

// Get retrieves a cached response
func (c *Cache) Get(command, errorMsg string) (*llm.Response, bool) {
	key := c.makeKey(command, errorMsg)
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	// Update hit counter and last used
	entry.Hits++
	entry.LastUsed = time.Now().Unix()

	// Parse the fix into a response
	response := &llm.Response{
		Cached: true,
	}

	// Parse the cached fix text
	lines := strings.Split(entry.Fix, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "✓"):
			response.Suggestion = strings.TrimSpace(strings.TrimPrefix(line, "✓"))
		case strings.HasPrefix(line, "Root:"):
			response.RootCause = strings.TrimSpace(strings.TrimPrefix(line, "Root:"))
		case strings.HasPrefix(line, "Tip:"):
			response.Tip = strings.TrimSpace(strings.TrimPrefix(line, "Tip:"))
		}
	}

	return response, true
}

// Set stores a response in the cache
func (c *Cache) Set(command, errorMsg string, response *llm.Response) error {
	key := c.makeKey(command, errorMsg)

	// Format the response as text
	fix := fmt.Sprintf("✓ %s\nRoot: %s", response.Suggestion, response.RootCause)
	if response.Tip != "" {
		fix += fmt.Sprintf("\nTip: %s", response.Tip)
	}

	c.entries[key] = &Entry{
		Command:   command,
		Error:     errorMsg,
		Fix:       fix,
		Timestamp: time.Now().Unix(),
		Hits:      0,
		LastUsed:  time.Now().Unix(),
	}

	return c.save()
}

// makeKey creates a cache key from command and error
func (c *Cache) makeKey(command, errorMsg string) string {
	// Use first line of error for key
	errorFirstLine := strings.Split(errorMsg, "\n")[0]
	data := command + "::" + errorFirstLine
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// load loads the cache from disk
func (c *Cache) load() error {
	data, err := os.ReadFile(c.file)
	if err != nil {
		return err
	}

	// Try to load with current format
	if err := json.Unmarshal(data, &c.entries); err != nil {
		// If it fails, might be old bash format or corrupted
		// Start with empty cache instead of failing
		c.entries = make(map[string]*Entry)
		return nil
	}

	return nil
}

// save saves the cache to disk
func (c *Cache) save() error {
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.file, data, 0644)
}

// Stats returns cache statistics
func (c *Cache) Stats() map[string]interface{} {
	totalHits := 0
	for _, entry := range c.entries {
		totalHits += entry.Hits
	}

	return map[string]interface{}{
		"total_entries": len(c.entries),
		"total_hits":    totalHits,
		"cache_file":    c.file,
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() error {
	c.entries = make(map[string]*Entry)
	return c.save()
}

