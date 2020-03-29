package account

import (
	"fmt"
	"math"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/cfanatic/go-expenses/database"
	"github.com/cfanatic/go-expenses/datasheet"
)

type Expense struct {
	Path string
	ds   *datasheet.Datasheet
	db   *database.Database
	ti   time.Time
	res  map[string]float64
	cnt  map[string]int
	tot  float64
	cat  []string
}

func (exp *Expense) Init() {
	exp.ds = datasheet.New(exp.Path, TAB)
	exp.db = database.New(ADDRESS, NAME, COLLECT)
	exp.res = map[string]float64{}
	exp.cnt = map[string]int{}
}

func (exp *Expense) Run() {
	exp.label()
	exp.analyze()
	exp.filter()
	exp.sort()
}

func (exp *Expense) Plot() {
}

func (exp *Expense) Print(filter ...string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 18, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n %s %s\n", fmt.Sprintf("%s", exp.ti.Month()), fmt.Sprintf("%d", exp.ti.Year()))
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "Label", "Euro / Month", "Percent / Month", "Count / Month")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "-----", "------------", "---------------", "-------------")
	var cats *[]string
	switch {
	case len(filter) != 0:
		cats = &filter
	default:
		cats = &exp.cat
	}
	for _, cat := range *cats {
		fmt.Fprintf(w, "\n %s\t%.2f\t%.f\t%d", cat, exp.res[cat], math.Round(exp.res[cat]/exp.tot*100), exp.cnt[cat])
	}
	fmt.Fprintf(w, "\n\n")
}

func (exp *Expense) Export() [][]string {
	var out [][]string
	var cnt int

	out = make([][]string, len(exp.cat)+3)
	for cat := range exp.cat {
		out[cat] = make([]string, 0)
	}

	out[0] = append(out[0], []string{
		fmt.Sprintf("%s %d", exp.ti.Month(), exp.ti.Year()),
		"Euro / Month",
		"Percent / Month",
		"Count / Month",
	}...)

	for idx, cat := range exp.cat {
		out[idx+1] = append(out[idx+1], []string{
			fmt.Sprintf("%s", cat),
			fmt.Sprintf("%.2f", exp.res[cat]),
			fmt.Sprintf("%.f", math.Round(exp.res[cat]/exp.tot*100)),
			fmt.Sprintf("%d", exp.cnt[cat]),
		}...)
		cnt = cnt + exp.cnt[cat]
	}

	out[len(exp.cat)+2] = append(out[len(exp.cat)+2], []string{
		"Total",
		fmt.Sprintf("%.0f €", exp.tot),
		"100 %",
		fmt.Sprintf("%d", cnt),
	}...)

	return out
}

func (exp *Expense) label() {
	var export []datasheet.Content
	var label string
	var err error

	if GUI == true {
		return
	}

	if export, err = exp.ds.Content(); err != nil {
		panic("Error during datasheet export!")
	}

	for _, trans := range export {
		if item, err := exp.db.Document("payee", trans.Payee); err == nil {
			fmt.Printf("| %s | %s | %.2f | -> %s\n\n", item.Payee, item.Desc, item.Amount, item.Label)
			date, _ := time.Parse("01-02-06", trans.Date)
			exp.db.Save(database.Content{
				Date:   date,
				Payee:  trans.Payee,
				Desc:   trans.Desc,
				Amount: trans.Amount,
				Label:  item.Label,
			})
		} else {
			fmt.Printf("\n%s\n%.2f\n%s\nLabel: ", trans.Payee, trans.Amount, trans.Desc)
			fmt.Scanln(&label)
			fmt.Println()
			date, _ := time.Parse("01-02-06", trans.Date)
			exp.db.Save(database.Content{
				Date:   date,
				Payee:  trans.Payee,
				Desc:   trans.Desc,
				Amount: trans.Amount,
				Label:  label,
			})
		}
	}
}

func (exp *Expense) analyze() {
	export, _ := exp.ds.Content()
	labels, _ := exp.db.Labels("label")

	dateIn := export[len(export)-1].Date
	dateOut := export[0].Date

	exp.ti, _ = time.Parse("01-02-06", dateIn)

	for _, label := range labels {
		content, _ := exp.db.Content("label", label.(string), dateIn, dateOut)
		exp.cnt[label.(string)] = len(content)
		for _, trans := range content {
			exp.res[trans.Label] = exp.res[trans.Label] + (-1.0 * float64(trans.Amount))
			exp.tot = exp.tot + (-1.0 * float64(trans.Amount))
		}
	}
}

func (exp *Expense) filter() {
	for _, label := range FILTER {
		exp.tot = exp.tot - exp.res[label]
		delete(exp.res, label)
	}
}

func (exp *Expense) sort() {
	exp.cat = make([]string, 0, len(exp.res))
	for cat := range exp.res {
		exp.cat = append(exp.cat, cat)
	}
	sort.Strings(exp.cat)
}
