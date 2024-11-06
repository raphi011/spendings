package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
)

type model struct {
	curView       string
	table         table.Model
	menu          list.Model
	spendingInput spendingInput
	spendings     []spending
}

func app(spendings []spending) *model {
	m := &model{spendings: spendings}

	renderSpendingsTable(m)
	renderMenu(m)
	renderSpendingInput(m)

	return m
}

func (m model) onNewSpending(s spending) {
	m.spendings = append(m.spendings, s)
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.curView = "menu"
			return m, cmd
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.curView == "" || m.curView == "menu" {
		m.menu, cmd = m.updateMenu(msg)
	} else if m.curView == "table" {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.spendingInput, cmd = m.spendingInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.curView == "" || m.curView == "menu" {
		return baseStyle.Render(m.menu.View()) + "\n"
	} else if m.curView == "table" {
		return baseStyle.Render(m.table.View()) + "\n"
	} else if m.curView == "input" {
		return m.spendingInput.View()
	}

	panic("unknown view")
}
