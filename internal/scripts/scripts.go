// Package scripts handles script discovery and management for go-pwr.
package scripts

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/list"
)

// Item represents a script or directory item.
type Item struct {
	name string
	path string
	tags *ScriptTags
}

// Title returns the display name of the item.
func (s Item) Title() string { return s.name }

// Description returns the path of the item.
func (s Item) Description() string { return s.path }

// FilterValue returns the value used for filtering (includes tags).
func (s Item) FilterValue() string { 
	filterValue := s.name
	if s.tags != nil {
		// Add tag values to filter string for better search
		for _, tag := range s.tags.Tags {
			filterValue += " " + tag.Value + " " + tag.Category
		}
	}
	return filterValue
}

// GetTags returns the tags for this item
func (s Item) GetTags() *ScriptTags {
	return s.tags
}

// IsDirectory returns true if the item is a directory.
func (s Item) IsDirectory() bool {
	return strings.HasSuffix(s.name, "/")
}

// IsScript returns true if the item is a supported script file.
func (s Item) IsScript() bool {
	ext := strings.ToLower(filepath.Ext(s.name))
	return ext == ".sh" || ext == ".ps1" || ext == ".bat" || ext == ".cmd"
}

// Cache provides thread-safe caching for script contents.
type Cache struct {
	mu    sync.RWMutex
	cache map[string]string
}

// NewCache creates a new script cache.
func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]string),
	}
}

// Get retrieves content from cache.
func (sc *Cache) Get(path string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	content, exists := sc.cache[path]
	return content, exists
}

// Set stores content in cache.
func (sc *Cache) Set(path, content string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache[path] = content
}

// Clear empties the cache.
func (sc *Cache) Clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache = make(map[string]string)
}

// GetItems returns all script items in the given directory.
func GetItems(root string) []list.Item {
	var items []list.Item
	entries, err := os.ReadDir(root)
	if err != nil {
		return items
	}

	// Pre-allocate slice with estimated capacity
	items = make([]list.Item, 0, len(entries))

	// Sort entries: directories first, then files, both alphabetically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}
		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files and directories
		if strings.HasPrefix(name, ".") {
			continue
		}

		path := filepath.Join(root, name)
		if entry.IsDir() {
			items = append(items, Item{name: name + "/", path: path, tags: nil})
		} else {
			// Only include supported script files
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".sh" || ext == ".ps1" || ext == ".bat" || ext == ".cmd" {
				// Parse tags for script files
				tags, _ := ParseTags(path) // Ignore errors, just use nil
				items = append(items, Item{name: name, path: path, tags: tags})
			}
		}
	}
	return items
}

// GetAllScriptsRecursively returns all script items recursively from the given directory
func GetAllScriptsRecursively(root string) []list.Item {
	var allItems []list.Item
	
	// Use a helper function to recursively walk directories
	var walkDir func(string, string)
	walkDir = func(currentPath, relativePath string) {
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return
		}
		
		for _, entry := range entries {
			name := entry.Name()
			// Skip hidden files and directories
			if strings.HasPrefix(name, ".") {
				continue
			}
			
			fullPath := filepath.Join(currentPath, name)
			displayPath := name
			if relativePath != "" {
				displayPath = filepath.Join(relativePath, name)
			}
			
			if entry.IsDir() {
				// Recursively walk subdirectories
				walkDir(fullPath, displayPath)
			} else {
				// Only include supported script files
				ext := strings.ToLower(filepath.Ext(name))
				if ext == ".sh" || ext == ".ps1" || ext == ".bat" || ext == ".cmd" {
					// Parse tags for script files
					tags, _ := ParseTags(fullPath) // Ignore errors, just use nil
					allItems = append(allItems, Item{
						name: displayPath, // Show relative path from root
						path: fullPath,    // Keep full path for execution
						tags: tags,
					})
				}
			}
		}
	}
	
	walkDir(root, "")
	
	// Sort items alphabetically by display name
	sort.Slice(allItems, func(i, j int) bool {
		return strings.ToLower(allItems[i].(Item).name) < strings.ToLower(allItems[j].(Item).name)
	})
	
	return allItems
}

// FilterItemsByTags filters script items by tags
func FilterItemsByTags(items []list.Item, searchTags []string) []list.Item {
	if len(searchTags) == 0 {
		return items
	}
	
	var filteredItems []list.Item
	for _, item := range items {
		if scriptItem, ok := item.(Item); ok {
			if scriptItem.IsDirectory() {
				// Always include directories
				filteredItems = append(filteredItems, item)
			} else if scriptItem.IsScript() && scriptItem.tags != nil {
				// Check if script matches ALL of the search tags (AND logic)
				if scriptItem.tags.HasAllTags(searchTags) {
					filteredItems = append(filteredItems, item)
				}
			}
		}
	}
	return filteredItems
}

// ReadContent reads the content of a script file, using cache when possible.
func ReadContent(path string, cache *Cache) string {
	// Check cache first
	if content, exists := cache.Get(path); exists {
		return content
	}

	data, err := os.ReadFile(path)
	if err != nil {
		errMsg := "Error reading file: " + err.Error()
		cache.Set(path, errMsg) // Cache errors too to avoid repeated attempts
		return errMsg
	}

	content := string(data)
	cache.Set(path, content)
	return content
}

// SanitizeContent removes or escapes characters that can break border rendering.
func SanitizeContent(content string) string {
	// Simple, reliable content cleaning
	content = strings.ReplaceAll(content, "\r\n", "\n") // Normalize line endings
	content = strings.ReplaceAll(content, "\r", "\n")   // Handle old Mac line endings
	content = strings.ReplaceAll(content, "\x00", "")   // Remove null bytes

	// Remove any stray ANSI sequences that might interfere
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Simple line length limit to prevent layout issues
		if len(line) > 200 { // Conservative limit
			lines[i] = line[:200] + "..."
		}
	}

	return strings.Join(lines, "\n")
}

// getBatCommand returns the appropriate bat command for the system.
func getBatCommand() string {
	// Check for bat first (most distributions)
	if _, err := exec.LookPath("bat"); err == nil {
		return "bat"
	}
	// Check for batcat (Ubuntu/Debian)
	if _, err := exec.LookPath("batcat"); err == nil {
		return "batcat"
	}
	return ""
}

// ReadContentWithHighlighting reads the content of a script file with syntax highlighting using bat.
func ReadContentWithHighlighting(path string, cache *Cache) string {
	// Check cache first
	if content, exists := cache.Get(path); exists {
		return content
	}

	// Use bat for syntax highlighting
	batCmd := getBatCommand()
	if batCmd == "" {
		// If bat is not available, show a helpful message
		content := "bat is not installed. Install it for syntax highlighting:\n" +
			"Windows: winget install sharkdp.bat\n" +
			"macOS: brew install bat\n" +
			"Ubuntu: sudo apt install bat\n\n" +
			"Raw file content:\n" + readRawContent(path)
		cache.Set(path, content)
		return content
	}

	// Use bat with DarkNeon theme
	cmd := exec.Command(batCmd, 
		"--theme=DarkNeon",
		"--color=always",
		"--style=numbers",
		"--paging=never",
		path,
	)
	
	output, err := cmd.Output()
	if err != nil {
		// If bat fails, show error and raw content
		content := "bat failed: " + err.Error() + "\n\nRaw content:\n" + readRawContent(path)
		cache.Set(path, content)
		return content
	}
	
	content := string(output)
	cache.Set(path, content)
	return content
}

// readRawContent reads file content without caching (helper function)
func readRawContent(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "Error reading file: " + err.Error()
	}
	return SanitizeContent(string(data))
}
