package dynreadme

import (
	"fmt"
	"os"
	"strings"

	"github.com/willabides/mdtable"
)

type Update struct{}

func (u *Update) updateContent(readmePath, markerText, mdText, table, tableOptions string) error {
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	startMarker := fmt.Sprintf("<!-- %s_START -->", markerText)
	endMarker := fmt.Sprintf("<!-- %s_END -->", markerText)

	var updatedReadmeContent string
	if table == "true" {
		var mdArray []string
		mdArray = strings.Split(mdText, ";")
		options := parseTableOptions(tableOptions)

		var data [][]string
		for _, md := range mdArray {
			data = append(data, strings.Split(md, ","))
		}

		updatedReadmeContent = replaceBetweenMarkers(string(readmeContent), startMarker, endMarker, fmt.Sprintf("%s\n%s\n%s", startMarker, string(mdtable.Generate(data, options...)), endMarker))
	} else if table == "false" {
		updatedReadmeContent = replaceBetweenMarkers(string(readmeContent), startMarker, endMarker, fmt.Sprintf("%s\n%s\n%s", startMarker, mdText, endMarker))
	}

	if err = os.WriteFile(readmePath, []byte(updatedReadmeContent), 0644); err != nil {
		return fmt.Errorf("error updating README: %w", err)
	}
	return nil
}

func parseTableOptions(tableOptions string) []mdtable.Option {
	var options []mdtable.Option
	parts := strings.Split(tableOptions, ",")

	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "align-"):
			alignment := parseAlignment(strings.TrimPrefix(part, "align-"))
			options = append(options, mdtable.Alignment(alignment))

		case strings.HasPrefix(part, "col-"):
			colParts := strings.Split(part, "-")
			if len(colParts) >= 3 {
				colNum := parseColumnNumber(colParts[1])
				switch colParts[2] {
				case "align":
					if len(colParts) == 4 {
						alignment := parseAlignment(colParts[3])
						options = append(options, mdtable.ColumnAlignment(colNum, alignment))
					}
				case "w":
					if len(colParts) == 4 {
						width := parseColumnWidth(colParts[3])
						options = append(options, mdtable.ColumnMinWidth(colNum, width))
					}
				}
			}
		}
	}
	return options
}

func parseAlignment(alignment string) mdtable.Align {
	switch alignment {
	case "left":
		return mdtable.AlignLeft
	case "right":
		return mdtable.AlignRight
	case "center":
		return mdtable.AlignCenter
	default:
		return mdtable.AlignDefault
	}
}

func parseColumnNumber(num string) int {
	var colNum int
	fmt.Sscanf(num, "%d", &colNum)
	return colNum
}

func parseColumnWidth(width string) int {
	var w int
	fmt.Sscanf(width, "%d", &w)
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
