package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/cfanatic/go-expense/account"
	"github.com/cfanatic/go-expense/database"
	"github.com/cfanatic/go-expense/datasheet"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	NORMAL = iota
	BOLD
	UNDERLINE
)

const (
	NONE = iota
	LABEL
	LOAD
)

var (
	FILTER = []string{""}
)

type Gui struct {
	widgets.QWidget

	qapp    *widgets.QApplication
	lapp    *widgets.QVBoxLayout
	lbutton *widgets.QHBoxLayout

	twidget     *widgets.QTabWidget
	twsettings  *widgets.QWidget
	twdata      *widgets.QWidget
	twmonth     *widgets.QWidget
	twyear      *widgets.QWidget
	twlsettings *widgets.QVBoxLayout
	twldata     *widgets.QHBoxLayout
	twlmonth    *widgets.QVBoxLayout
	twlyear     *widgets.QHBoxLayout

	flist *widgets.QListWidget
	tview *widgets.QTableView
	tlist *gui.QStandardItemModel

	ds    *datasheet.Datasheet
	db    *database.Database
	dlist []database.Content

	sfile string
	spath string
	smode int

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
	w.twsettings = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twdata = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twmonth = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twyear = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twlsettings = widgets.NewQVBoxLayout()
	w.twldata = widgets.NewQHBoxLayout()
	w.twlmonth = widgets.NewQVBoxLayout()
	w.twlyear = widgets.NewQHBoxLayout()

	w.twidget.AddTab(w.twsettings, "Settings")
	w.twidget.AddTab(w.twdata, "Data")
	w.twidget.AddTab(w.twmonth, "Month")
	w.twidget.AddTab(w.twyear, "Year")
	w.twidget.SetTabEnabled(1, false)
	w.twidget.SetTabEnabled(2, false)
	w.twidget.SetTabEnabled(3, false)

	w.twsettings.SetLayout(w.twlsettings)
	w.twmonth.SetLayout(w.twlmonth)
	w.twyear.SetLayout(w.twlyear)

	w.flist = widgets.NewQListWidget(nil)
	w.twlsettings.AddWidget(w.flist, 0, 0)

	w.tview = widgets.NewQTableView(nil)
	w.tlist = gui.NewQStandardItemModel(nil)

	w.tview.SetModel(w.tlist)
	w.tview.SetContextMenuPolicy(core.Qt__CustomContextMenu)
	w.twldata.AddWidget(w.tview, 0, 0)
	w.twdata.SetLayout(w.twldata)

	blabel := widgets.NewQPushButton2("Label", nil)
	bload := widgets.NewQPushButton2("Load", nil)
	bquit := widgets.NewQPushButton2("Quit", nil)

	w.lbutton.AddWidget(blabel, 0, 0)
	w.lbutton.AddWidget(bload, 0, 0)
	w.lbutton.AddWidget(bquit, 0, 0)
	w.lapp.AddWidget(w.twidget, 0, 0)
	w.lapp.AddLayout(w.lbutton, 0)
	w.SetLayout(w.lapp)

	w.twidget.ConnectTabBarClicked(w.analyze)
	w.flist.ConnectItemDoubleClicked(w.removefilter)
	w.tview.ConnectCustomContextMenuRequested(w.addfilter)
	w.tlist.ConnectItemChanged(w.update)
	blabel.ConnectClicked(func(bool) { w.data(LABEL) })
	bload.ConnectClicked(func(bool) { w.data(LOAD) })
	bquit.ConnectClicked(func(bool) { w.qapp.Exit(0) })

	w.ConnectKeyPressEvent(w.keypressevent)
}

func (w *Gui) data(mode int) {
	var export []datasheet.Content
	var err error

	w.smode = mode

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
		tmp := strings.Split(f[0], "/")
		w.sfile = tmp[len(tmp)-1]
		w.spath = core.NewQFileInfo3(f[0]).AbsoluteDir().AbsolutePath()
	} else {
		return
	}

	if export, err = w.ds.Content(); err == nil {
		w.twidget.SetTabEnabled(1, true)
		w.twidget.SetTabEnabled(2, true)
		w.twidget.SetTabEnabled(3, true)
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

	if w.tlist.RowCount(core.NewQModelIndex()) > 0 {
		w.flist.Clear()
		w.tlist.RemoveRows(0, w.tlist.RowCount(core.NewQModelIndex()), core.NewQModelIndex())
		w.dlist = []database.Content{}
	}

	for _, trans := range export {
		items := []*gui.QStandardItem{}
		id := w.db.Hash(trans)
		if item, err := w.db.Document("datasheet", id); err == nil {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(item.Label),
			)
		} else if item, err := w.db.Document("payee", trans.Payee); err == nil {
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
		switch w.smode {
		case LABEL:
			items[0].SetEditable(false)
			items[0].SetSelectable(false)
		case LOAD:
			for i := 0; i < 4; i++ {
				items[i].SetEditable(false)
				items[i].SetSelectable(false)
			}
		}
		items[1].SetToolTip(trans.Desc)
		w.tlist.AppendRow(items)
		w.dlist = append(w.dlist, database.Content{
			Date:      date,
			Payee:     trans.Payee,
			Desc:      trans.Desc,
			Amount:    trans.Amount,
			Datasheet: id,
		})
	}

	width := float32(w.Geometry().Width())
	w.tlist.SetHorizontalHeaderLabels([]string{"Date", "Payee", "Amount", "Label"})
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Interactive)
	w.tview.HorizontalHeader().SetStretchLastSection(true)
	w.tview.HorizontalHeader().ResizeSection(0, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(1, int(width/2))
	w.tview.HorizontalHeader().ResizeSection(2, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(3, int(width/5))
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)

	w.save()
}

func (w *Gui) save() {
	for row := 0; row < w.tlist.RowCount(core.NewQModelIndex()); row++ {
		trans := []string{}
		for col := 0; col < w.tlist.ColumnCount(core.NewQModelIndex()); col++ {
			index := w.tlist.Index(row, col, core.NewQModelIndex())
			data := w.tlist.Data(index, int(core.Qt__DisplayRole))
			trans = append(trans, data.ToString())
		}
		trans = append(trans, w.tlist.Item(row, 1).ToolTip(), w.dlist[row].Datasheet)
		if doc := w.document(trans); len(doc.Label) > 0 {
			w.db.Save(doc)
		}
	}
}

func (w *Gui) update(item *gui.QStandardItem) {
	var trans []string

	trans = make([]string, 0)
	for col := 0; col < w.tlist.ColumnCount(core.NewQModelIndex()); col++ {
		index := w.tlist.Index(item.Row(), col, core.NewQModelIndex())
		data := w.tlist.Data(index, int(core.Qt__DisplayRole))
		trans = append(trans, data.ToString())
	}
	trans = append(trans, w.tlist.Item(item.Row(), 1).ToolTip(), "id_not_known_yet")

	dsold := datasheet.Content{
		Date:   w.dlist[item.Row()].Date.Format("01-02-06"),
		Payee:  w.dlist[item.Row()].Payee,
		Desc:   w.dlist[item.Row()].Desc,
		Amount: w.dlist[item.Row()].Amount,
	}
	dsnew := datasheet.Content{
		Date:   w.document(trans).Date.Format("01-02-06"),
		Payee:  w.document(trans).Payee,
		Desc:   w.document(trans).Desc,
		Amount: w.document(trans).Amount,
	}

	trans[5] = w.db.Hash(dsnew)
	dbold := w.dlist[item.Row()]
	dbnew := w.document(trans)

	w.ds.Update(dsold, dsnew)
	w.db.Update(dbold, dbnew)

	w.dlist[item.Row()] = w.document(trans)
}

func (w *Gui) addfilter(position *core.QPoint) {
	menu := widgets.NewQMenu(nil)
	index := w.tview.IndexAt(position)
	item := w.tlist.ItemFromIndex(index)

	if len(item.Text()) > 0 {
		action := widgets.NewQAction2(fmt.Sprintf("Exclude: %s", item.Text()), w)
		action.ConnectTriggered(func(bool) {
			list := w.flist.FindItems(item.Text(), core.Qt__MatchExactly)
			if len(list) == 0 {
				row := widgets.NewQListWidgetItem(w.flist, 0)
				font := row.Font()
				font.SetPointSize(font.PointSize() + 0)
				row.SetFont(font)
				row.SetText(item.Text())
			}
		})

		menu.AddActions([]*widgets.QAction{action})

		if item.Column()+1 == 4 {
			menu.Popup(w.tview.Viewport().MapToGlobal(position), action)
		}
	}
}

func (w *Gui) removefilter(item *widgets.QListWidgetItem) {
	w.flist.TakeItem(w.flist.Row(item))
	w.flist.RemoveItemWidget(item)
}

func (w *Gui) analyze(index int) {
	switch w.twidget.TabText(index) {
	case "Month":
		w.month()
	case "Year":
		w.year()
	}
}

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qapp.Exit(0)
	}
}

func (w *Gui) InitWith(qapp *widgets.QApplication) {
	w.qapp = qapp
}
