package view

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
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
	go a.logs().logs.Clear()
	go a.Logs("docs", "server1")
	go a.Logs("docs", "server2")

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

func (a *App) testCounter() *tview.TextView {
	tv := tview.NewTextView()
	tv.SetBackgroundColor(tcell.ColorDefault)
	tv.SetTextColor(tcell.ColorLightGoldenrodYellow)
	tv.SetScrollable(false)
	tv.SetChangedFunc(func() {
		a.Draw()
	})
	go func() {
		var n int
		for {
			n++
			fmt.Fprintf(tv, "%d ", n)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return tv
}

func (a *App) nodeStatus() *tview.Table {
	t := tview.NewTable()
	// t.SetBackgroundColor(tcell.ColorDefault)
	// t.SetTextColor(tcell.ColorLightGoldenrodYellow)
	// t.SetSelectionChangedFunc(func() {
	// 	a.Draw()
	// })

	// go func() {
	// 	// t.setCo
	// }()

	return t
}

func (a *App) Watch() {
	ticker := time.NewTicker(a.refreshRate)
	for {
		select {
		case <-ticker.C:
			info, err := a.client.NomadInfo()
			if err != nil {
				panic(err)
			}
			a.nomadInfo().InfoUpdated(info)
		}
	}
}

func (a *App) Logs(job, task string) {
	reader, err := a.client.Logs(job, task)
	if err != nil {
		panic(err)
	}

	lw := &LogWriter{
		app:    a,
		name:   task,
		writer: tview.ANSIWriter(a.logs().logs, "#d8dee9", "#2e3440"),
		color:  colorFor(task),
	}

	defer reader.Close()
	_, err = io.Copy(lw, reader)
	if err != nil {
		panic(err)
	}
}

type LogWriter struct {
	app    *App
	name   string
	color  int
	writer io.Writer
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	lw.app.QueueUpdateDraw(func() {
		fmt.Fprint(lw.writer, lw.format(p))
	})
	return len(p), nil
}

const colorFmt = "\033[38;5;%dm%s\033[0m"

func colorize(s string, c int) string {
	return fmt.Sprintf(colorFmt, c, s)
}

func (lw *LogWriter) format(p []byte) string {
	service := colorize(lw.name, lw.color)
	return fmt.Sprintf("%s: %s", service, string(p))
}

func colorFor(n string) int {
	return rand.Intn(250)
}
