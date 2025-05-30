package entimport

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FormatGeneratedFiles formats all Go files in the specified directory
func FormatGeneratedFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		formatted := improveCodeFormatting(string(content))
		if formatted != string(content) {
			err = os.WriteFile(path, []byte(formatted), info.Mode())
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// improveCodeFormatting improves the formatting of generated code
func improveCodeFormatting(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// Check if this line contains Fields() method with inline field definitions
		if strings.Contains(line, "return []ent.Field{") && strings.Contains(line, "field.") {
			// Find the indentation
			indent := ""
			for i, char := range line {
				if char != '\t' && char != ' ' {
					indent = line[:i]
					break
				}
			}

			// Extract fields content
			start := strings.Index(line, "return []ent.Field{")
			end := strings.LastIndex(line, "}")
			if start != -1 && end > start {
				fieldsStart := start + len("return []ent.Field{")
				fieldsContent := line[fieldsStart:end]

				// Split fields by ", field." pattern
				var fields []string
				if strings.TrimSpace(fieldsContent) != "" {
					parts := strings.Split(fieldsContent, ", field.")
					for i, part := range parts {
						part = strings.TrimSpace(part)
						if part == "" {
							continue
						}
						if i == 0 {
							// First part already has "field." prefix
							fields = append(fields, part)
						} else {
							// Add "field." prefix back
							fields = append(fields, "field."+part)
						}
					}
				}

				// Fix JSON field parameter order and reconstruct with proper formatting
				result = append(result, indent+"return []ent.Field{")
				for _, field := range fields {
					field = strings.TrimSpace(field)
					if field != "" {
						// Fix JSON field parameter order
						field = fixJSONFieldOrder(field)
						result = append(result, indent+"\t"+field+",")
					}
				}
				result = append(result, indent+"}")
				continue
			}
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// fixJSONFieldOrder fixes the parameter order for JSON fields
func fixJSONFieldOrder(fieldStr string) string {
	// Pattern: field.JSON("name").Optional().Comment("comment", struct{}{})
	// Should be: field.JSON("name", struct{}{}).Optional().Comment("comment")
	if strings.Contains(fieldStr, "field.JSON(") && strings.Contains(fieldStr, "Comment(") && strings.Contains(fieldStr, "struct{}{}") {
		// Find the JSON field name
		startName := strings.Index(fieldStr, "field.JSON(\"") + len("field.JSON(\"")
		endName := strings.Index(fieldStr[startName:], "\"") + startName
		if startName < len("field.JSON(\"") || endName <= startName {
			return fieldStr
		}
		fieldName := fieldStr[startName:endName]

		// Find the comment content
		commentPattern := "Comment(\""
		commentStart := strings.Index(fieldStr, commentPattern)
		if commentStart == -1 {
			return fieldStr
		}
		commentContentStart := commentStart + len(commentPattern)
		commentContentEnd := strings.Index(fieldStr[commentContentStart:], "\"") + commentContentStart
		if commentContentEnd <= commentContentStart {
			return fieldStr
		}
		comment := fieldStr[commentContentStart:commentContentEnd]

		// Check if Optional() is present
		optional := ""
		if strings.Contains(fieldStr, ").Optional().") {
			optional = ".Optional()"
		}

		// Reconstruct the field with correct parameter order
		return fmt.Sprintf("field.JSON(\"%s\", struct{}{})%s.Comment(\"%s\")", fieldName, optional, comment)
	}
	return fieldStr
}
