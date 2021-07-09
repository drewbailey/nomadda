package view

import (
	"io"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

const (
	initialLogMsg = "Waiting for logs...\n"
)

type Log struct {
	*tview.Flex
	app *App

	writer io.Writer
	logs   *Details
}

func NewLog(app *App) *Log {
	return &Log{
		Flex: tview.NewFlex(),
		app:  app,
	}
}

func (l *Log) Init() {
	l.SetBorder(true)
	l.SetBorderColor(tcell.ColorLightSlateGray)
	l.SetBackgroundColor(tcell.ColorDefault)
	l.SetDirection(tview.FlexRow)

	l.logs = NewDetails(l.app)
	l.logs.Init()
	l.logs.SetBorderPadding(0, 0, 1, 1)
	l.logs.SetText(initialLogMsg)
	l.logs.SetWrap(true)

	// l.writer = tview.ANSIWriter(l.logs, "white", "black")
	l.AddItem(l.logs, 0, 10, true)
}
