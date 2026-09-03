package tskcli

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// formatProjectListTable renders human-readable project list rows as an
// aligned table (NAME ORIGIN LOCATION [TASKS]). Empty cells are blank. tasks is
// right-aligned when present. Caller prints the footer separately.
func formatProjectListTable(rows []projectListJSONRow, includeTasks bool) string {
	if len(rows) == 0 {
		return ""
	}
	headers := []string{"NAME", "ORIGIN", "LOCATION"}
	if includeTasks {
		headers = append(headers, "TASKS")
	}
	table := make([][]string, 0, 1+len(rows))
	table = append(table, headers)
	for _, r := range rows {
		row := []string{r.Name, r.Origin, r.Location}
		if includeTasks {
			row = append(row, fmt.Sprintf("%d", r.Tasks))
		}
		table = append(table, row)
	}
	return formatAlignedTable(table, includeTasks)
}

// formatAlignedTable pads columns to max rune width with a 2-space gap.
// When rightAlignLast is true, the last column is padLeft (for TASKS).
func formatAlignedTable(rows [][]string, rightAlignLast bool) string {
	if len(rows) == 0 {
		return ""
	}
	cols := len(rows[0])
	widths := make([]int, cols)
	for _, row := range rows {
		for j := 0; j < cols && j < len(row); j++ {
			if w := runeWidth(row[j]); w > widths[j] {
				widths[j] = w
			}
		}
	}
	var b strings.Builder
	for _, row := range rows {
		for j := 0; j < cols; j++ {
			cell := ""
			if j < len(row) {
				cell = row[j]
			}
			if rightAlignLast && j == cols-1 {
				b.WriteString(padLeft(cell, widths[j]))
			} else {
				b.WriteString(padRight(cell, widths[j]))
			}
			if j+1 < cols {
				b.WriteString("  ")
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func runeWidth(s string) int {
	return utf8.RuneCountInString(s)
}

func padRight(s string, width int) string {
	w := runeWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func padLeft(s string, width int) string {
	w := runeWidth(s)
	if w >= width {
		return s
	}
	return strings.Repeat(" ", width-w) + s
}
