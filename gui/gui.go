package gui

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Gui struct {
	widgets.QWidget

	qApp     *widgets.QApplication
	hlayout  *widgets.QHBoxLayout
	vlayout  *widgets.QVBoxLayout
	tview    *widgets.QTableView
	tviewm   *gui.QStandardItemModel
	bconnect *widgets.QPushButton
	bload    *widgets.QPushButton
	bsave    *widgets.QPushButton
	bprint   *widgets.QPushButton

	_ func() `constructor:"init"`

	_ func() `slot:"connect"`
}

func (w *Gui) init() {
	geometry := widgets.QApplication_Desktop().AvailableGeometry(0)
	w.SetWindowTitle("go-expensegui")
	w.SetGeometry2(0, 0, 600, 600)
	w.Move2((geometry.Width()-w.Width())/2, (geometry.Height()-w.Height())/2)

	w.ConnectKeyPressEvent(w.keypressevent)

	w.hlayout = widgets.NewQHBoxLayout()
	w.vlayout = widgets.NewQVBoxLayout()
	w.tview = widgets.NewQTableView(nil)
	w.tviewm = gui.NewQStandardItemModel(nil)
	w.bconnect = widgets.NewQPushButton2("Connect", nil)
	w.bload = widgets.NewQPushButton2("Load", nil)
	w.bsave = widgets.NewQPushButton2("Save", nil)
	w.bprint = widgets.NewQPushButton2("Print", nil)

	w.tview.SetModel(w.tviewm)
	w.hlayout.AddWidget(w.bconnect, 0, 0)
	w.hlayout.AddWidget(w.bload, 0, 0)
	w.hlayout.AddWidget(w.bsave, 0, 0)
	w.hlayout.AddWidget(w.bprint, 0, 0)
	w.vlayout.AddWidget(w.tview, 0, 0)
	w.vlayout.AddLayout(w.hlayout, 0)

	w.SetLayout(w.vlayout)

	w.bconnect.ConnectClicked(w.connect)
}

func (w *Gui) connect(bool) {
	header := []string{"Date", "Payee", "Amount", "Label"}
	items := []*gui.QStandardItem{}
	// item1 := gui.NewQStandardItem2("Test1")
	// item2 := gui.NewQStandardItem2("Test2")
	// item3 := gui.NewQStandardItem2("Test3")
	// item4 := gui.NewQStandardItem2("Test4")
	// items = append(items, item1, item2, item3, item4)
	w.tviewm.SetHorizontalHeaderLabels(header)
	w.tviewm.AppendRow(items)
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)
}

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qApp.Exit(0)
	}
}

func (w *Gui) InitWith(qApp *widgets.QApplication) {
	w.qApp = qApp
}
