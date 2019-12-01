package gui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cfanatic/go-expense/account"
	"github.com/cfanatic/go-expense/database"
	"github.com/cfanatic/go-expense/datasheet"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Gui struct {
	widgets.QWidget

	qapp    *widgets.QApplication
	lapp    *widgets.QVBoxLayout
	lbutton *widgets.QHBoxLayout

	twidget  *widgets.QTabWidget
	twlabel  *widgets.QWidget
	twmonth  *widgets.QWidget
	twyear   *widgets.QWidget
	twinfo   *widgets.QWidget
	twlload  *widgets.QHBoxLayout
	twlmonth *widgets.QHBoxLayout
	twlyear  *widgets.QHBoxLayout
	twldebug *widgets.QHBoxLayout

	vlayout *widgets.QVBoxLayout
	tview   *widgets.QTableView
	list    *gui.QStandardItemModel

	ds    *datasheet.Datasheet
	db    *database.Database
	listd []database.Content

	_ func() `constructor:"init"`

	_ func() `slot:"connect"`
}

func (w *Gui) init() {
	geometry := widgets.QApplication_Desktop().AvailableGeometry(0)
	w.SetWindowTitle("go-expensegui")
	w.SetGeometry2(0, 0, int(float32(geometry.Width())/2.2), geometry.Height()/2)
	w.Move2((geometry.Width()-w.Width())/2, (geometry.Height()-w.Height())/2)

	w.lapp = widgets.NewQVBoxLayout()
	w.lbutton = widgets.NewQHBoxLayout()

	w.twidget = widgets.NewQTabWidget(nil)
	w.twlabel = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twmonth = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twyear = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twinfo = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twlload = widgets.NewQHBoxLayout()
	w.twlmonth = widgets.NewQHBoxLayout()
	w.twlyear = widgets.NewQHBoxLayout()
	w.twldebug = widgets.NewQHBoxLayout()

	w.vlayout = widgets.NewQVBoxLayout()
	w.tview = widgets.NewQTableView(nil)
	w.list = gui.NewQStandardItemModel(nil)

	w.tview.SetModel(w.list)
	w.vlayout.AddWidget(w.tview, 0, 0)

	w.twlabel.SetLayout(w.vlayout)
	w.twidget.AddTab(w.twinfo, "Info")
	w.twidget.AddTab(w.twlabel, "Label")
	w.twidget.AddTab(w.twmonth, "Month")
	w.twidget.AddTab(w.twyear, "Year")
	w.twidget.SetTabEnabled(1, false)
	w.twidget.SetTabEnabled(2, false)
	w.twidget.SetTabEnabled(3, false)

	blabel := widgets.NewQPushButton2("Label", nil)
	bquit := widgets.NewQPushButton2("Quit", nil)

	w.lbutton.AddWidget(blabel, 0, 0)
	w.lbutton.AddWidget(bquit, 0, 0)
	w.lapp.AddWidget(w.twidget, 0, 0)
	w.lapp.AddLayout(w.lbutton, 0)
	w.SetLayout(w.lapp)

	blabel.ConnectClicked(w.label)
	bquit.ConnectClicked(func(bool) { w.qapp.Exit(0) })
	w.list.ConnectItemChanged(w.update)
	w.ConnectKeyPressEvent(w.keypressevent)
}

func (w *Gui) update(item *gui.QStandardItem) {
	var trans []string

	trans = make([]string, 0)
	for col := 0; col < w.list.ColumnCount(core.NewQModelIndex()); col++ {
		index := w.list.Index(item.Row(), col, core.NewQModelIndex())
		data := w.list.Data(index, int(core.Qt__DisplayRole))
		trans = append(trans, data.ToString())
	}

	dsold := datasheet.Content{
		Date:   w.listd[item.Row()].Date.Format("01-02-06"),
		Payee:  w.listd[item.Row()].Payee,
		Amount: w.listd[item.Row()].Amount,
	}
	dsnew := datasheet.Content{
		Date:   w.document(trans).Date.Format("01-02-06"),
		Payee:  w.document(trans).Payee,
		Amount: w.document(trans).Amount,
	}
	dbold := w.listd[item.Row()]
	dbnew := w.document(trans)

	w.ds.Update(dsold, dsnew)
	w.db.Update(dbold, dbnew)

	w.listd[item.Row()] = w.document(trans)
}

func (w *Gui) label(bool) {
	var export []datasheet.Content
	var err error

	if f := widgets.QFileDialog_GetOpenFileNames(nil, "Open datasheet", core.QDir_HomePath(), "*.xlsx", "", 0); len(f) > 0 {
		w.ds = datasheet.New(f[0], account.TAB)
		w.db = database.New(account.ADDRESS, account.NAME, account.COLLECT)
		if w.ds.Err != nil {
			widgets.QMessageBox_Critical(nil,
				"Cannot open datasheet",
				"Error during datasheet initialization attempt!",
				widgets.QMessageBox__Default,
				widgets.QMessageBox__Default,
			)
			return
		}
		if w.db.Err != nil {
			widgets.QMessageBox_Critical(nil,
				"Cannot open database",
				"Error during database initialization attempt!",
				widgets.QMessageBox__Default,
				widgets.QMessageBox__Default,
			)
			return
		}
	} else {
		return
	}

	if export, err = w.ds.Content(); err == nil {
		w.twidget.SetTabEnabled(1, true)
		w.twidget.SetCurrentIndex(1)
	} else {
		widgets.QMessageBox_Critical(nil,
			"Cannot import datasheet",
			"Error during datasheet import!",
			widgets.QMessageBox__Default,
			widgets.QMessageBox__Default,
		)
		return
	}

	if w.list.RowCount(core.NewQModelIndex()) > 0 {
		w.list.RemoveRows(0, w.list.RowCount(core.NewQModelIndex()), core.NewQModelIndex())
		w.listd = []database.Content{}
	}

	for _, trans := range export {
		items := []*gui.QStandardItem{}
		if item, err := w.db.Document("payee", trans.Payee); err == nil {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(item.Label),
			)
		} else {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(""),
			)
		}
		date, _ := time.Parse("01-02-06", trans.Date)
		items[0].SetEditable(false)
		items[0].SetSelectable(false)
		items[1].SetToolTip(trans.Description)
		w.list.AppendRow(items)
		w.listd = append(w.listd, database.Content{Date: date, Payee: trans.Payee, Amount: trans.Amount})
	}

	width := float32(w.Geometry().Width())
	w.list.SetHorizontalHeaderLabels([]string{"Date", "Payee", "Amount", "Label"})
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Interactive)
	w.tview.HorizontalHeader().SetStretchLastSection(true)
	w.tview.HorizontalHeader().ResizeSection(0, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(1, int(width/2))
	w.tview.HorizontalHeader().ResizeSection(2, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(3, int(width/5))
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)

	w.save(true)
}

func (w *Gui) save(bool) {
	for row := 0; row < w.list.RowCount(core.NewQModelIndex()); row++ {
		trans := []string{}
		for col := 0; col < w.list.ColumnCount(core.NewQModelIndex()); col++ {
			index := w.list.Index(row, col, core.NewQModelIndex())
			data := w.list.Data(index, int(core.Qt__DisplayRole))
			trans = append(trans, data.ToString())
		}
		if doc := w.document(trans); len(doc.Label) > 0 {
			w.db.Save(doc)
		}
	}
}

func (w *Gui) document(trans []string) database.Content {
	date, _ := time.Parse("01-02-06", trans[0])
	payee := trans[1]
	amount, _ := strconv.ParseFloat(trans[2], 32)
	label := trans[3]
	return database.Content{Date: date, Payee: payee, Amount: float32(amount), Label: label}
}

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qapp.Exit(0)
	}
}

func (w *Gui) InitWith(qapp *widgets.QApplication) {
	w.qapp = qapp
}
