package gui

import (
	"io/ioutil"
	"strings"

	"github.com/cfanatic/go-expense/account"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func (w *Gui) month() {
	var acc account.IAccount
	var res [][]string

	row := func(style int, items ...string) *widgets.QWidget {
		widget := widgets.NewQWidget(nil, 0)
		layout := widgets.NewQHBoxLayout()
		for _, item := range items {
			label := widgets.NewQLabel2(item, nil, core.Qt__Widget)
			switch font := label.Font(); style {
			case BOLD:
				font.SetBold(true)
				label.SetFont(font)
			case UNDERLINE:
				font.SetUnderline(true)
				label.SetFont(font)
			default:
			}
			layout.AddWidget(label, 0, 0)
		}
		widget.SetLayout(layout)
		return widget
	}

	path := w.spath + "/" + w.sfile
	acc = &account.Expense{Path: path}
	acc.Init()
	acc.Run()

	res = acc.Export()

	sarea := widgets.NewQScrollArea(nil)
	swidget := widgets.NewQWidget(nil, 0)
	slayout := widgets.NewQVBoxLayout2(swidget)

	slayout.AddWidget(row(BOLD, res[0]...), 0, 0)

	for idx := 1; idx < len(res)-2; idx++ {
		slayout.AddWidget(row(NORMAL, res[idx]...), 0, 0)
	}

	slayout.AddWidget(row(UNDERLINE, res[len(res)-1]...), 0, 0)

	spacer := widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding)
	slayout.AddSpacerItem(spacer)

	sarea.SetWidget(swidget)
	sarea.SetWidgetResizable(true)

	if w.twlmonth.Count() > 0 {
		tmp := w.twlmonth.ItemAt(0).Widget()
		tmp.Hide()
		w.twlmonth.RemoveWidget(tmp)
	}
	w.twlmonth.InsertWidget(0, sarea, 0, 0)
}

func (w *Gui) year() {
	var acc account.IAccount
	var exps []*account.Expense
	var res [][]string

	row := func(style int, items ...string) *widgets.QWidget {
		widget := widgets.NewQWidget(nil, 0)
		layout := widgets.NewQHBoxLayout()
		for idx, item := range items {
			label := widgets.NewQLabel2(item, nil, core.Qt__Widget)
			switch font := label.Font(); style {
			case BOLD:
				font.SetBold(true)
				label.SetFont(font)
			case UNDERLINE:
				font.SetUnderline(true)
				label.SetFont(font)
			default:
			}
			if idx == 1 || idx == len(items)-5 {
				spacer := widgets.NewQSpacerItem(60, 0, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
				layout.AddSpacerItem(spacer)
			}
			layout.AddWidget(label, 0, 0)
		}
		widget.SetLayout(layout)
		return widget
	}

	if dir, _ := ioutil.ReadDir(w.spath); len(dir) > 0 {
		for _, file := range dir {
			str := strings.Split(file.Name(), ".")
			if strings.Contains(str[0], "~$") || str[1] != "xlsx" {
				continue
			}
			path := w.spath + "/" + file.Name()
			acc = &account.Expense{Path: path}
			acc.Init()
			acc.Run()
			exps = append(exps, acc.(*account.Expense))
		}
		acc = &account.Expenses{Exp: exps}
		acc.Init()
		acc.Run()
	}

	res = acc.Export()

	sarea := widgets.NewQScrollArea(nil)
	swidget := widgets.NewQWidget(nil, 0)
	slayout := widgets.NewQVBoxLayout2(swidget)

	slayout.AddWidget(row(BOLD, res[0]...), 0, 0)

	for idx := 1; idx < len(res); idx++ {
		slayout.AddWidget(row(NORMAL, res[idx]...), 0, 0)
	}

	spacer := widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding)
	slayout.AddSpacerItem(spacer)

	sarea.SetWidget(swidget)
	sarea.SetWidgetResizable(true)

	if w.twlyear.Count() > 0 {
		tmp := w.twlyear.ItemAt(0).Widget()
		tmp.Hide()
		w.twlyear.RemoveWidget(tmp)
	}
	w.twlyear.InsertWidget(0, sarea, 0, 0)
}
