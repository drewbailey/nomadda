package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var nomadLogo = []string{
	`@@@@@@@@@@@@@@@*@@@@@@@@@@@@@@`,
	`@@@@@@@@@@@*********@@@@@@@@@@`,
	`@@@@@@@*****************@@@@@@`,
	`@@@******************* ******@`,
	`(((***************     *****,,`,
	`(((((((***********     *,,,,,,`,
	`(((((((((        *     ,,,,,,,`,
	`((((((((               ,,,,,,,`,
	`((((((((     (((     ,,,,,,,,,`,
	`((((((((     (((,,,,,,,,,,,,,,`,
	`((((((((     (((,,,,,,,,,,,,,,`,
	`@@@((((( (((((((,,,,,,,,,,,,@@`,
	`@@@@@@@(((((((((,,,,,,,,@@@@@@`,
	`@@@@@@@@@@@@((((,,,,@@@@@@@@@@`,
}

type Logo struct {
	*tview.Flex
	logo *tview.TextView
}

func NewLogo() *Logo {
	logo := &Logo{
		logo: logo(),
		Flex: tview.NewFlex(),
	}

	logo.SetBackgroundColor(tcell.ColorDefault)
	logo.logo.SetBackgroundColor(tcell.ColorDefault)
	logo.SetDirection(tview.FlexRow)
	logo.AddItem(logo.logo, 120, 1, false)
	logo.refresh()
	return logo
}

func (l *Logo) refresh() {
	l.logo.Clear()
	for i, s := range nomadLogo {
		fmt.Fprintf(l.logo, "[%s::b]%s", "slategrey", s)
		if i+1 < len(nomadLogo) {
			fmt.Fprintf(l.logo, "\n")
		}
	}
}

func logo() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWordWrap(false)
	v.SetTextAlign(tview.AlignCenter)
	v.SetDynamicColors(true)
	return v
}
