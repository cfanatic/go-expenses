package gui

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/cfanatic/go-expenses/account"
	"github.com/cfanatic/go-expenses/database"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func (w *Gui) reset() {
	account.GUI = true
	account.FILTER = []string{}
	for idx := 0; idx < w.flist.Count(); idx++ {
		label := w.flist.Item(idx)
		account.FILTER = append(account.FILTER, label.Text())
	}
}

func (w *Gui) month() {
	var acc account.IAccount
	var res [][]string

	w.reset()

	acc = &account.Expense{Path: w.spath + "/" + w.sfile}
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
	slayout.AddStretch(1)

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

	w.reset()

	if dir, _ := ioutil.ReadDir(w.spath); len(dir) > 0 {
		for _, file := range dir {
			str := strings.Split(file.Name(), ".")
			if len(str) == 1 || strings.Contains(str[0], "~$") || str[1] != "xlsx" {
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

	for idx := 1; idx < len(res)-2; idx++ {
		slayout.AddWidget(row(NORMAL, res[idx]...), 0, 0)
	}

	slayout.AddWidget(row(UNDERLINE, res[len(res)-1]...), 0, 0)

	spacer := widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding)
	slayout.AddSpacerItem(spacer)
	slayout.AddStretch(1)

	sarea.SetWidget(swidget)
	sarea.SetWidgetResizable(true)

	if w.twlyear.Count() > 0 {
		tmp := w.twlyear.ItemAt(0).Widget()
		tmp.Hide()
		w.twlyear.RemoveWidget(tmp)
	}
	w.twlyear.InsertWidget(0, sarea, 0, 0)
}

func (w *Gui) document(trans []string) database.Content {
	date, _ := time.Parse("01-02-06", trans[0])
	amount, _ := strconv.ParseFloat(trans[2], 32)
	return database.Content{
		Date:      date,
		Payee:     trans[1],
		Desc:      trans[4],
		Amount:    float32(amount),
		Label:     trans[3],
		Datasheet: trans[5],
	}
}

func (w *Gui) integrity() {
	dateIn := w.dlist[len(w.dlist)-1].Date
	dateOut := w.dlist[0].Date

	remote, _ := w.db.Content("", "", dateIn.Format("01-02-06"), dateOut.Format("01-02-06"))
	local := []string{}

	for row := 0; row < w.tlist.RowCount(core.NewQModelIndex()); row++ {
		index := w.tlist.Index(row, 3, core.NewQModelIndex())
		data := w.tlist.Data(index, int(core.Qt__DisplayRole))
		if len(data.ToString()) > 0 {
			local = append(local, data.ToString())
		}
	}

	if len(remote) == len(local) {
		w.lwarning.SetVisible(false)
	} else {
		w.lwarning.SetText(fmt.Sprintf("REDUNDANT DATA: R-%d L-%d", len(remote), len(local)))
		w.lwarning.SetStyleSheet("font-weight: bold; color: red")
		w.lwarning.SetVisible(true)
	}
}

func row(style int, items ...string) *widgets.QWidget {
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
