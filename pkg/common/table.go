package common

import (
	"github.com/olekukonko/tablewriter"
)

// TableType exports table formatting.
type TableType int

// Tabler exports table formatting const object.
var Tabler = &tablerList{
	CSV:      csv,
	Markdown: markdown,
	Preview:  preview,
	Basic:    basic,
	Vertical: vertical,
}

const (
	csv TableType = iota + 1
	markdown
	preview
	basic
	vertical
)

type tablerList struct {
	CSV      TableType
	Markdown TableType
	Preview  TableType
	Basic    TableType
	Vertical TableType
}

// String converts TableType to its value.
func (ttype TableType) String() string {
	return [...]string{"CSV", "Markdown", "Preview", "Basic"}[ttype]
}

// SetStyleForTable configures table style based on type.
func (ttype TableType) SetStyleForTable(table *tablewriter.Table, size int) {
	switch ttype {
	case preview:
		headerStyles := make([]tablewriter.Colors, size)
		columnStyles := make([]tablewriter.Colors, size)
		for i := 0; i < size; i++ {
			headerStyles[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
		}
		// columnStyles[0] = tablewriter.Colors{tablewriter.Bold}
		table.SetHeaderColor(headerStyles...)
		table.SetColumnColor(columnStyles...)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	case csv:
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetColumnSeparator(",")
		table.SetHeaderLine(false)
	case markdown:
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	case basic:
		headerStyles := make([]tablewriter.Colors, size)
		for i := 0; i < size; i++ {
			headerStyles[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
		}
		table.SetHeaderColor(headerStyles...)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator(" ")
		table.SetColumnSeparator(" ")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
	case vertical:
		headerStyles := make([]tablewriter.Colors, 2)
		columnStyles := make([]tablewriter.Colors, 2)
		headerStyles[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
		columnStyles[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
		table.SetHeaderColor(headerStyles...)
		table.SetColumnColor(columnStyles...)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoFormatHeaders(false)
		table.SetCenterSeparator(" ")
		table.SetColumnSeparator(" ")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
	}
}
