package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	Main  *Pages
	views map[string]tview.Primitive
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		Main:        NewPages(),
		views: map[string]tview.Primitive{
			"logo": NewLogo(),
		},
	}
	// clear out to allow for default color
	app.Application.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})
	return app
}

func (a *App) Init() {
	a.Main.SetBackgroundColor(tcell.ColorDefault)
	a.Main.AddPage("header", a.header(), true, true)
	a.SetRoot(a.Main, true)
}

func (a *App) header() tview.Primitive {
	header := tview.NewFlex()
	header.SetDirection(tview.FlexColumn)
	header.SetBackgroundColor(tcell.ColorDefault)
	header.AddItem(a.Logo(), 120, 1, false)
	return header
}

func (a *App) Logo() tview.Primitive {
	return a.views["logo"]
}

type Pages struct {
	*tview.Pages
}

func NewPages() *Pages {
	return &Pages{
		Pages: tview.NewPages(),
	}
}
