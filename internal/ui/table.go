package ui

import "github.com/rivo/tview"

type Table struct {
	*tview.Table

	app        *App
	selectedFn func(string) string
}
