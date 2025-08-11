package scripts

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Tag represents a tag with its category and value
type Tag struct {
	Category string
	Value    string
}

// ScriptTags holds all tags for a script
type ScriptTags struct {
	Path string
	Tags []Tag
}

// TagCategory represents different types of tags
type TagCategory struct {
	Name   string
	Values []string
}

// ParseTags extracts tags from a script file
func ParseTags(filePath string) (*ScriptTags, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tags []Tag
	scanner := bufio.NewScanner(file)
	inTagsSection := false
	
	// Regex patterns for different tag formats
	tagsStartPattern := regexp.MustCompile(`^#\*Tags:?\s*$`)
	categoryPattern := regexp.MustCompile(`^#\s*([A-Za-z_]+):\s*(.+)$`)
	oldFormatPattern := regexp.MustCompile(`^#([a-zA-Z_]+)\s+#([a-zA-Z_]+)`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Stop parsing after first 50 lines to avoid parsing entire file
		if len(tags) > 0 && !inTagsSection {
			break
		}
		
		// Check if we're starting the tags section
		if tagsStartPattern.MatchString(line) {
			inTagsSection = true
			continue
		}
		
		// If in tags section, parse category-based tags
		if inTagsSection {
			if matches := categoryPattern.FindStringSubmatch(line); matches != nil {
				category := strings.TrimSpace(matches[1])
				values := strings.Fields(matches[2])
				for _, value := range values {
					value = strings.Trim(value, ",")
					if value != "" {
						tags = append(tags, Tag{
							Category: strings.ToLower(category),
							Value:    strings.ToLower(value),
						})
					}
				}
			} else if line != "" && !strings.HasPrefix(line, "#") {
				// End of tags section if we hit non-comment content
				break
			}
		} else {
			// Parse old format tags (for backward compatibility)
			if matches := oldFormatPattern.FindAllStringSubmatch(line, -1); matches != nil {
				for _, match := range matches {
					tag := strings.ToLower(strings.TrimPrefix(match[0], "#"))
					if tag != "" {
						tags = append(tags, Tag{
							Category: "legacy",
							Value:    tag,
						})
					}
				}
			}
		}
		
		// Stop after empty line following tags
		if inTagsSection && line == "" {
			break
		}
	}
	
	return &ScriptTags{
		Path: filePath,
		Tags: tags,
	}, scanner.Err()
}

// HasTag checks if a script has a specific tag
func (st *ScriptTags) HasTag(category, value string) bool {
	for _, tag := range st.Tags {
		if strings.EqualFold(tag.Category, category) && strings.EqualFold(tag.Value, value) {
			return true
		}
	}
	return false
}

// HasAnyTag checks if a script has any tag from a list (OR logic)
func (st *ScriptTags) HasAnyTag(searchTags []string) bool {
	for _, searchTag := range searchTags {
		searchTag = strings.ToLower(strings.TrimSpace(searchTag))
		for _, tag := range st.Tags {
			if strings.Contains(tag.Value, searchTag) || strings.Contains(tag.Category, searchTag) {
				return true
			}
		}
	}
	return false
}

// HasAllTags checks if a script has ALL tags from a list (AND logic)
func (st *ScriptTags) HasAllTags(searchTags []string) bool {
	if len(searchTags) == 0 {
		return true
	}
	
	// Check that every search tag is found in the script's tags
	for _, searchTag := range searchTags {
		searchTag = strings.ToLower(strings.TrimSpace(searchTag))
		found := false
		
		// Look for this search tag in any of the script's tags
		for _, tag := range st.Tags {
			if strings.Contains(tag.Value, searchTag) || strings.Contains(tag.Category, searchTag) {
				found = true
				break
			}
		}
		
		// If any search tag is not found, return false
		if !found {
			return false
		}
	}
	
	// All search tags were found
	return true
}

// GetTagsByCategory returns all tags for a specific category
func (st *ScriptTags) GetTagsByCategory(category string) []string {
	var values []string
	for _, tag := range st.Tags {
		if strings.EqualFold(tag.Category, category) {
			values = append(values, tag.Value)
		}
	}
	return values
}

// GetAllCategories returns all unique categories
func (st *ScriptTags) GetAllCategories() []string {
	categorySet := make(map[string]bool)
	for _, tag := range st.Tags {
		categorySet[tag.Category] = true
	}
	
	categories := make([]string, 0, len(categorySet))
	for category := range categorySet {
		categories = append(categories, category)
	}
	return categories
}

// SearchScriptsByTags searches for scripts that match given tags
func SearchScriptsByTags(scriptPaths []string, searchTags []string) ([]*ScriptTags, error) {
	var matchingScripts []*ScriptTags
	
	for _, path := range scriptPaths {
		scriptTags, err := ParseTags(path)
		if err != nil {
			continue // Skip files with errors
		}
		
		if len(searchTags) == 0 || scriptTags.HasAllTags(searchTags) {
			matchingScripts = append(matchingScripts, scriptTags)
		}
	}
	
	return matchingScripts, nil
}

// GetAllTagsFromDirectory recursively finds all unique tags in a directory
func GetAllTagsFromDirectory(rootPath string) (map[string][]string, error) {
	allTags := make(map[string][]string)
	
	items := GetItems(rootPath)
	for _, item := range items {
		if scriptItem, ok := item.(Item); ok && scriptItem.IsScript() {
			scriptTags, err := ParseTags(scriptItem.Description())
			if err != nil {
				continue
			}
			
			for _, tag := range scriptTags.Tags {
				if _, exists := allTags[tag.Category]; !exists {
					allTags[tag.Category] = []string{}
				}
				
				// Add tag value if not already present
				found := false
				for _, existingValue := range allTags[tag.Category] {
					if strings.EqualFold(existingValue, tag.Value) {
						found = true
						break
					}
				}
				if !found {
					allTags[tag.Category] = append(allTags[tag.Category], tag.Value)
				}
			}
		}
	}
	
	return allTags, nil
}
