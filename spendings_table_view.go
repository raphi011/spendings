package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func formatAmount(amountInCents int) string {
	return fmt.Sprintf("%.2f €", float32(amountInCents)/100)
}

func renderSpendingsTable(m *model) {
	columns := []table.Column{
		{Title: "Date", Width: 10},
		{Title: "Reference", Width: 40},
		{Title: "Labels", Width: 20},
		{Title: "Amount", Width: 10},
	}

	rows := make([]table.Row, len(m.spendings))

	for i, s := range m.spendings {
		rows[i] = table.Row{s.date.Format("2006-01-02"), s.reference, strings.Join(s.labels, ","), formatAmount(s.amountInCents)}
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	m.table.SetStyles(s)

}
