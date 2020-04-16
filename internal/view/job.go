package view

import "github.com/derailed/tview"

type Job struct {
	*tview.Table
}

func NewJob() *Job {
	return &Job{
		Table: tview.NewTable(),
	}
}
