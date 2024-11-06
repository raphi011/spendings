package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	dt = iota
	amt
	ref
	lbl
)

func renderSpendingInput(m *model) {
	var inputs []textinput.Model = make([]textinput.Model, 4)

	inputs[dt] = textinput.New()
	inputs[dt].Placeholder = "YYYY-MM-DD"
	inputs[dt].Focus()
	inputs[dt].CharLimit = 10
	inputs[dt].Width = 10
	inputs[dt].Prompt = ""
	inputs[dt].Validate = dateValidator

	inputs[amt] = textinput.New()
	inputs[amt].Placeholder = "13,55"
	inputs[amt].CharLimit = 10
	inputs[amt].Width = 10
	inputs[amt].Prompt = ""
	inputs[amt].Validate = amountValidator

	inputs[ref] = textinput.New()
	inputs[ref].Placeholder = "Grocery shopping"
	inputs[ref].CharLimit = 30
	inputs[ref].Width = 10
	inputs[ref].Prompt = ""

	inputs[lbl] = textinput.New()
	inputs[lbl].Placeholder = "Music"
	inputs[lbl].CharLimit = 10
	inputs[lbl].Width = 10
	inputs[lbl].Prompt = ""

	m.spendingInput = spendingInput{
		inputs:        inputs,
		focused:       0,
		onNewSpending: m.onNewSpending,
	}
}

var dateChars = regexp.MustCompile("[0-9-]{0,10}")
var amountChars = regexp.MustCompile("^[0-9]+,?([0-9]{1,2})?$")

func amountValidator(s string) error {
	if !amountChars.MatchString(s) {
		// slog.Info(fmt.Sprintf("%s does not match", s))
		return errors.New("amount contains invalid chars")
	}
	// slog.Info(fmt.Sprintf("%s matches", s))

	return nil
}

func dateValidator(s string) error {
	if !dateChars.MatchString(s) {
		return errors.New("date contains invalid chars")
	}

	if len(s) == 10 {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
	}

	return nil
}

type spendingInput struct {
	inputs        []textinput.Model
	focused       int
	onNewSpending func(s spending)
	err           error
}

func (m spendingInput) View() string {
	const (
		hotPink  = lipgloss.Color("#FF06B7")
		darkGray = lipgloss.Color("#767676")
	)

	var (
		inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
		continueStyle = lipgloss.NewStyle().Foreground(darkGray)
	)
	return fmt.Sprintf(
		`
 %s
 %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Date"),
		m.inputs[dt].View(),

		inputStyle.Width(30).Render("Amount"),
		m.inputs[amt].View(),

		inputStyle.Width(30).Render("Reference"),
		m.inputs[ref].View(),

		inputStyle.Width(30).Render("Label"),
		m.inputs[lbl].View(),

		continueStyle.Render("Add ->"),
	) + "\n"
}

func (m spendingInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m spendingInput) Update(msg tea.Msg) (spendingInput, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			s, err := m.newSpending()
			if err != nil {
				m.err = err
				panic(err)
			} else {
				m.onNewSpending(s)
			}
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m spendingInput) newSpending() (spending, error) {
	date, err := time.Parse("2006-01-02", m.inputs[dt].Value())
	if err != nil {
		return spending{}, fmt.Errorf("invalid date: %w", err)
	}

	amount, err := strconv.Atoi(m.inputs[amt].Value())
	if err != nil {
		return spending{}, fmt.Errorf("invalid amount: %w", err)
	}

	reference := m.inputs[ref].Value()
	label := m.inputs[lbl].Value()

	return spending{
		date:          date,
		amountInCents: amount,
		labels:        []string{label},
		reference:     reference,
		payee:         "",
	}, nil
}

// nextInput focuses the next input field
func (m *spendingInput) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *spendingInput) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
