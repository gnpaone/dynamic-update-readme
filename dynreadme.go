package dynreadme

import (
	"fmt"
	"log"
	"os"
	"strings"

	table "github.com/gnpaone/dynamic-update-readme/helpers"
)

// UpdateContent parses and updates content between markers
func UpdateContent(readmePath, markerText, mdText, isTable, tableOptions string) error {
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	startMarker := fmt.Sprintf("<!-- %s_START -->", markerText)
	endMarker := fmt.Sprintf("<!-- %s_END -->", markerText)

	var updatedReadmeContent string
	if isTable == "true" {
		// var mdArray []string
		var mdArray = strings.Split(mdText, ";")
		options := parseTableOptions(tableOptions)

		var data [][]string
		for _, md := range mdArray {
			data = append(data, strings.Split(md, ","))
		}

		updatedReadmeContent = replaceBetweenMarkers(string(readmeContent), startMarker, endMarker, fmt.Sprintf("%s\n%s\n%s", startMarker, string(table.Generate(data, options...)), endMarker))
	} else if isTable == "false" {
		updatedReadmeContent = replaceBetweenMarkers(string(readmeContent), startMarker, endMarker, fmt.Sprintf("%s\n%s\n%s", startMarker, mdText, endMarker))
	}

	if err = os.WriteFile(readmePath, []byte(updatedReadmeContent), 0644); err != nil {
		return fmt.Errorf("error updating README: %w", err)
	}
	return nil
}

func parseTableOptions(tableOptions string) []table.Option {
	var options []table.Option
	parts := strings.Split(tableOptions, ",")

	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "align-"):
			alignment := parseAlignment(strings.TrimPrefix(part, "align-"))
			options = append(options, table.Alignment(alignment))

		case strings.HasPrefix(part, "col-"):
			colParts := strings.Split(part, "-")
			if len(colParts) >= 3 {
				colNum := parseColumnNumber(colParts[1])
				switch colParts[2] {
				case "align":
					if len(colParts) == 4 {
						alignment := parseAlignment(colParts[3])
						options = append(options, table.ColumnAlignment(colNum, alignment))
					}
				case "w":
					if len(colParts) == 4 {
						width := parseColumnWidth(colParts[3])
						options = append(options, table.ColumnMinWidth(colNum, width))
					}
				}
			}

		case strings.HasPrefix(part, "colH-"):
			colHParts := strings.Split(part, "-")
			if len(colHParts) >= 3 {
				colHNum := parseColumnNumber(colHParts[1])
				switch colHParts[2] {
				case "align":
					if len(colHParts) == 4 {
						alignment := parseAlignment(colHParts[3])
						options = append(options, table.ColumnHeaderAlignment(colHNum, alignment))
					}
				}
			}

		case strings.HasPrefix(part, "colT-"):
			colTParts := strings.Split(part, "-")
			if len(colTParts) >= 3 {
				colTNum := parseColumnNumber(colTParts[1])
				switch colTParts[2] {
				case "align":
					if len(colTParts) == 4 {
						alignment := parseAlignment(colTParts[3])
						options = append(options, table.ColumnTextAlignment(colTNum, alignment))
					}
				}
			}

		case strings.HasPrefix(part, "head-"):
			headParts := strings.Split(part, "-")
			if len(headParts) >= 2 {
				switch headParts[1] {
				case "align":
					if len(headParts) == 3 {
						alignment := parseAlignment(headParts[2])
						options = append(options, table.HeaderAlignment(alignment))
					}
				}

			}

		case strings.HasPrefix(part, "text-"):
			textParts := strings.Split(part, "-")
			if len(textParts) >= 2 {
				switch textParts[1] {
				case "align":
					if len(textParts) == 3 {
						alignment := parseAlignment(textParts[2])
						options = append(options, table.TextAlignment(alignment))
					}
				}

			}
		}
	}
	return options
}

func parseAlignment(alignment string) table.Align {
	switch alignment {
	case "left":
		return table.AlignLeft
	case "right":
		return table.AlignRight
	case "center":
		return table.AlignCenter
	default:
		return table.AlignDefault
	}
}

func parseColumnNumber(num string) int {
	var colNum int
	if _, err := fmt.Sscanf(num, "%d", &colNum); err != nil {
		log.Fatalf("Error parsing column number: %s", err)
    	}
	return colNum
}

func parseColumnWidth(width string) int {
	var w int
	if _, err := fmt.Sscanf(width, "%d", &w); err != nil {
		log.Fatalf("Error parsing column width: %s", err)
    	}
	return w
}

func replaceBetweenMarkers(content, startMarker, endMarker, replacement string) string {
	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker) + len(endMarker)
	if startIdx == -1 || endIdx == -1 {
		return content
	}
	return content[:startIdx] + replacement + content[endIdx:]
}
