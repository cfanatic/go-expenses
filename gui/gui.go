package gui

import (
	"fmt"
	"log"
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

	qApp    *widgets.QApplication
	hlayout *widgets.QHBoxLayout
	vlayout *widgets.QVBoxLayout
	tview   *widgets.QTableView
	bload   *widgets.QPushButton
	bsave   *widgets.QPushButton
	bprint  *widgets.QPushButton
	bquit   *widgets.QPushButton
	list    *gui.QStandardItemModel
	listd   []database.Content

	ds *datasheet.Datasheet
	db *database.Database

	_ func() `constructor:"init"`

	_ func() `slot:"connect"`
}

func (w *Gui) init() {
	geometry := widgets.QApplication_Desktop().AvailableGeometry(0)
	w.SetWindowTitle("go-expensegui")
	w.SetGeometry2(0, 0, 1000, 700)
	w.Move2((geometry.Width()-w.Width())/2, (geometry.Height()-w.Height())/2)

	w.ConnectKeyPressEvent(w.keypressevent)

	w.hlayout = widgets.NewQHBoxLayout()
	w.vlayout = widgets.NewQVBoxLayout()
	w.tview = widgets.NewQTableView(nil)
	w.bload = widgets.NewQPushButton2("Load", nil)
	w.bsave = widgets.NewQPushButton2("Save", nil)
	w.bprint = widgets.NewQPushButton2("Print", nil)
	w.bquit = widgets.NewQPushButton2("Quit", nil)
	w.list = gui.NewQStandardItemModel(nil)
	w.listd = []database.Content{}

	w.bsave.SetEnabled(false)
	w.bprint.SetEnabled(false)

	w.tview.SetModel(w.list)
	w.hlayout.AddWidget(w.bload, 0, 0)
	w.hlayout.AddWidget(w.bsave, 0, 0)
	w.hlayout.AddWidget(w.bprint, 0, 0)
	w.hlayout.AddWidget(w.bquit, 0, 0)
	w.vlayout.AddWidget(w.tview, 0, 0)
	w.vlayout.AddLayout(w.hlayout, 0)

	w.SetLayout(w.vlayout)

	w.bload.ConnectClicked(w.load)
	w.bsave.ConnectClicked(w.save)
	w.bquit.ConnectClicked(func(bool) { w.qApp.Exit(0) })
	w.list.ConnectItemChanged(w.update)
	w.tview.HorizontalHeader().ConnectSectionResized(
		func(idx, old, new int) { log.Printf("Index: %d, Size: %d\n", idx, new) },
	)
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

func (w *Gui) load(bool) {
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
				"Error during database connection attempt!",
				widgets.QMessageBox__Default,
				widgets.QMessageBox__Default,
			)
			return
		}
	} else {
		return
	}

	if export, err = w.ds.Content(); err == nil {
		w.bsave.SetEnabled(true)
		w.bprint.SetEnabled(true)
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

	w.list.SetHorizontalHeaderLabels([]string{"Date", "Payee", "Amount", "Label"})
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Interactive)
	w.tview.HorizontalHeader().SetStretchLastSection(true)
	w.tview.HorizontalHeader().ResizeSection(0, 100)
	w.tview.HorizontalHeader().ResizeSection(1, 500)
	w.tview.HorizontalHeader().ResizeSection(2, 100)
	w.tview.HorizontalHeader().ResizeSection(3, 200)
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
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

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qApp.Exit(0)
	}
}

func (w *Gui) InitWith(qApp *widgets.QApplication) {
	w.qApp = qApp
}

func (w *Gui) document(trans []string) database.Content {
	date, _ := time.Parse("01-02-06", trans[0])
	amount, _ := strconv.ParseFloat(trans[2], 32)
	payee, label := trans[1], trans[3]
	return database.Content{Date: date, Payee: payee, Amount: float32(amount), Label: label}
}
