package view

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
	"github.com/hashicorp/nomadda/internal/structs"
)

type NomadInfo struct {
	*tview.Table
	app *App
}

func NewNomadInfo(a *App) *NomadInfo {
	return &NomadInfo{
		Table: tview.NewTable(),
		app:   a,
	}
}

func (n *NomadInfo) Init() {
	n.SetBorderPadding(0, 0, 1, 1)
	n.SetBackgroundColor(tcell.ColorDefault)
	n.layout()
}

func (n *NomadInfo) layout() {
	for i, s := range []string{"servers", "clients"} {
		n.SetCell(i, 0, n.sectionCell(s))
		n.SetCell(i, 1, n.infoCell("<loading>"))
	}
}

func (n *NomadInfo) sectionCell(s string) *tview.TableCell {
	cell := tview.NewTableCell(fmt.Sprintf("%s:", s))
	cell.SetAlign(tview.AlignLeft)
	return cell
}

func (n *NomadInfo) infoCell(s string) *tview.TableCell {
	cell := tview.NewTableCell(s)
	cell.SetExpansion(2)
	return cell
}

func (n *NomadInfo) setCell(row int, s string) int {
	n.GetCell(row, 1).SetText(s)
	return row + 1
}

func (n *NomadInfo) InfoUpdated(info structs.NomadInfo) {
	n.InfoChanged(info)
}

func (n *NomadInfo) InfoChanged(info structs.NomadInfo) {
	n.app.QueueUpdateDraw(func() {
		n.Clear()
		n.layout()
		row := n.setCell(0, info.Servers.GoString())
		row = n.setCell(row, info.Clients.GoString())
	})
}
