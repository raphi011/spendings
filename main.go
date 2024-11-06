package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) != 1 {
		slog.Error("Please pass a spending file")
		return
	}

	spendings, err := loadFile(files[0])
	if err != nil {
		slog.Error("Error loading spendings file", "err", err)
	}

	m := app(spendings)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type spending struct {
	date          time.Time
	amountInCents int
	labels        []string
	reference     string
	payee         string
}

func loadFile(path string) ([]spending, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse csv file: %w", err)
	}

	spendings := make([]spending, len(data))

	for i, row := range data {
		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, fmt.Errorf("date %q has invalid format: %w", row[0], err)
		}

		amount, err := strconv.Atoi(row[5])
		if err != nil {
			return nil, fmt.Errorf("amount %q is not a valid number: %w", row[5], err)
		}

		labels := strings.Split(row[2], ",")
		reference := row[3]
		payee := row[1]

		spendings[i] = spending{
			date:          date,
			amountInCents: amount,
			labels:        labels,
			reference:     reference,
			payee:         payee,
		}
	}

	return spendings, nil
}
