package view

import (
	"io"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/derailed/tview"
	"github.com/hashicorp/nomadda/internal/nomad"
	"github.com/hashicorp/nomadda/internal/ui"
)

type App struct {
	*ui.App
	client      *nomad.Client
	refreshRate time.Duration
}

func NewApp() *App {
	c, err := nomad.NewClient(&nomad.Config{})
	if err != nil {
		panic(err)
	}

	a := &App{
		App:         ui.NewApp(),
		client:      c,
		refreshRate: 1 * time.Second,
	}

	a.Views()["nomadInfo"] = NewNomadInfo(a)
	a.Views()["logs"] = NewLog(a)
	return a
}

func (a *App) Init() error {
	// init ui app
	a.App.Init()

	a.nomadInfo().Init()

	a.logs().Init()
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	main.AddItem(a.buildHeader(), 0, 1, false)
	main.AddItem(a.logs(), 0, 3, true)
	// main.AddItemAtIndex(0, a.buildHeader(), 15, 1, false)
	// main.AddItemAtIndex(1, a.logs(), 10, 5, false)
	a.Main.AddPage("main", main, true, false)

	// TODO why setroot here does somethings vs ui/app only like k9s
	a.App.SetRoot(main, true)

	// start polling for changes
	go a.Watch()
	go a.Logs()

	return nil
}

func (a *App) buildHeader() tview.Primitive {
	header := tview.NewFlex()
	header.SetDirection(tview.FlexColumn)
	header.AddItem(a.nomadInfo(), 0, 1, false)
	header.AddItem(a.Logo(), 0, 1, false)

	return header
}

func (a *App) nomadInfo() *NomadInfo {
	return a.Views()["nomadInfo"].(*NomadInfo)
}

func (a *App) logs() *Log {
	return a.Views()["logs"].(*Log)
}

func (a *App) Watch() {
	ticker := time.NewTicker(a.refreshRate)
	for {
		select {
		case <-ticker.C:
			info, err := a.client.NomadInfo()
			if err != nil {
				spew.Dump(err)
			}
			a.nomadInfo().InfoUpdated(info)
		}
	}
}

func (a *App) Logs() {
	reader, err := a.client.Logs()
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	_, err = io.Copy(a.logs(), reader)
	if err != nil {
		panic(err)
	}
}
