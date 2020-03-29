package account

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/gonum/stat"
	"github.com/ryanuber/columnize"
	"github.com/wcharczuk/go-chart"
)

type Expenses struct {
	Exp []*Expense
	res map[string]float64
	cnt map[string]int
	tot float64
	cat []string
}

func (exps *Expenses) Init() {
	exps.cnt = map[string]int{}
	exps.res = map[string]float64{}
	exps.tot = 0.0
}

func (exps *Expenses) Run() {
	exps.analyze()
	exps.filter()
	exps.sort()
}

func (exps *Expenses) Plot() {
	var title string
	var bars []chart.Value

	timeIn := fmt.Sprintf("%s %d",
		exps.Exp[0].ti.Month(),
		exps.Exp[0].ti.Year(),
	)
	timeOut := fmt.Sprintf("%s %d",
		exps.Exp[len(exps.Exp)-1].ti.Month(),
		exps.Exp[len(exps.Exp)-1].ti.Year(),
	)

	if timeIn == timeOut {
		title = fmt.Sprintf("%s", timeIn)
	} else {
		title = fmt.Sprintf("%s  —  %s", timeIn, timeOut)
	}

	for _, cat := range exps.cat {
		bar := chart.Value{Value: exps.res[cat], Label: cat}
		bars = append(bars, bar)
	}

	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Title:      title,
		Width:      1300,
		Height:     500,
		Bars:       bars,
		BarSpacing: 100,
		BarWidth:   20,
	}

	file, _ := os.Create("res/expenses/output.png")
	defer file.Close()
	graph.Render(chart.PNG, file)
}

func (exps *Expenses) Print(filter ...string) {
	var row strings.Builder
	var cats *[]string
	var out []string

	yearIn := fmt.Sprintf("%d", exps.Exp[0].ti.Year())
	yearOut := fmt.Sprintf("%d", exps.Exp[len(exps.Exp)-1].ti.Year())

	if yearIn == yearOut {
		row.WriteString(fmt.Sprintf("%s ||", yearIn))
	} else {
		row.WriteString(fmt.Sprintf("%s — %s ||", yearIn, yearOut))
	}

	for _, exp := range exps.Exp {
		row.WriteString(fmt.Sprintf("%s | ", exp.ti.Month()))
	}
	row.WriteString("|| TOTAL | AVERAGE | STDDEV | PERCENT | COUNT")
	out = append(out, row.String())
	out = append(out, "")

	switch {
	case len(filter) != 0:
		cats = &filter
	default:
		cats = &exps.cat
	}
	for _, cat := range *cats {
		row.Reset()
		row.WriteString(fmt.Sprintf("%s ||", cat))
		amt := []float64{}
		for _, exp := range exps.Exp {
			if val, ok := exp.res[cat]; ok {
				row.WriteString(fmt.Sprintf("%.2f | ", val))
				amt = append(amt, val)
			} else {
				row.WriteString(fmt.Sprintf("%.2f | ", 0.0))
				amt = append(amt, 0.0)
			}
		}
		row.WriteString(fmt.Sprintf("|| %.f | %.f | %.f | %.f | %d",
			math.Round(exps.res[cat]),
			math.Round(exps.res[cat]/float64(len(exps.Exp))),
			stat.StdDev(amt, nil),
			math.Round(exps.res[cat]/exps.tot*100),
			exps.cnt[cat]),
		)
		out = append(out, row.String())
	}

	row.Reset()
	row.WriteString("Total |")
	for _, exp := range exps.Exp {
		row.WriteString(fmt.Sprintf("| %.f", exp.tot))
	}
	out = append(out, "")
	out = append(out, row.String())

	config := columnize.DefaultConfig()
	config.Glue = "      "
	result := columnize.Format(out, config)
	fmt.Println(result)
}

func (exps *Expenses) Export() [][]string {
	var out [][]string
	var amt []float64
	var tot float64
	var cnt int

	out = make([][]string, len(exps.cat)+3)
	for cat := range exps.cat {
		out[cat] = make([]string, 0)
	}

	cnt = 0
	yearIn := fmt.Sprintf("%d", exps.Exp[0].ti.Year())
	yearOut := fmt.Sprintf("%d", exps.Exp[len(exps.Exp)-1].ti.Year())

	if yearIn == yearOut {
		out[0] = append(out[0], fmt.Sprintf("%s", yearIn))
	} else {
		out[0] = append(out[0], fmt.Sprintf("%s — %s", yearIn, yearOut))
	}
	for _, exp := range exps.Exp {
		out[0] = append(out[0], fmt.Sprintf("%s", exp.ti.Month()))
	}
	out[0] = append(out[0], []string{"TOTAL", "AVERAGE", "STDDEV", "PERCENT", "COUNT"}...)

	for idx, cat := range exps.cat {
		out[idx+1] = append(out[idx+1], fmt.Sprintf("%s", cat))
		amt = []float64{}
		for _, exp := range exps.Exp {
			if val, ok := exp.res[cat]; ok {
				out[idx+1] = append(out[idx+1], fmt.Sprintf("%.2f", val))
				amt = append(amt, val)
			} else {
				out[idx+1] = append(out[idx+1], fmt.Sprintf("%.2f", 0.0))
				amt = append(amt, 0.0)
			}
		}
		out[idx+1] = append(out[idx+1], []string{
			fmt.Sprintf("%.f", math.Round(exps.res[cat])),
			fmt.Sprintf("%.f", math.Round(exps.res[cat]/float64(len(exps.Exp)))),
			fmt.Sprintf("%.f", stat.StdDev(amt, nil)),
			fmt.Sprintf("%.f", math.Round(exps.res[cat]/exps.tot*100)),
			fmt.Sprintf("%d", exps.cnt[cat]),
		}...)
		cnt = cnt + exps.cnt[cat]
	}

	out[len(exps.cat)+2] = append(out[len(exps.cat)+2], "Total")
	amt = []float64{}
	tot = 0.0
	for _, exp := range exps.Exp {
		out[len(exps.cat)+2] = append(out[len(exps.cat)+2], fmt.Sprintf("%.f €", exp.tot))
		amt = append(amt, exp.tot)
		tot = tot + exp.tot
	}
	out[len(exps.cat)+2] = append(out[len(exps.cat)+2], []string{
		fmt.Sprintf("%.f €", tot),
		fmt.Sprintf("%.f €", tot/float64(len(exps.Exp))),
		fmt.Sprintf("%.f", stat.StdDev(amt, nil)),
		"100 %",
		fmt.Sprintf("%d", cnt),
	}...)

	return out
}

func (exps *Expenses) analyze() {
	for _, exp := range exps.Exp {
		for _, label := range exp.cat {
			exps.cnt[label] = exps.cnt[label] + exp.cnt[label]
			exps.res[label] = exps.res[label] + exp.res[label]
			exps.tot = exps.tot + exp.res[label]
		}
	}
}

func (exps *Expenses) filter() {
	for _, label := range FILTER {
		exps.tot = exps.tot - exps.res[label]
		delete(exps.res, label)
	}
}

func (exps *Expenses) sort() {
	exps.cat = make([]string, 0, len(exps.res))
	for cat := range exps.res {
		exps.cat = append(exps.cat, cat)
	}
	sort.Strings(exps.cat)
}
