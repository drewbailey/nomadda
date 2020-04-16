package ui

import "github.com/derailed/tview"

type Pages struct {
	*tview.Pages
}

func NewPages() *Pages {
	return &Pages{
		Pages: tview.NewPages(),
	}
}
