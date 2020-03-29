package account

const (
	// Excel tab name
	TAB = "expenses"
	// MongoDB database address
	ADDRESS = "mongodb://127.0.0.1:27017"
	// MongoDB database name
	NAME = "expenses"
	// MongoDB collection name
	COLLECT = "transaction"
)

var (
	FILTER = []string{""}
	GUI    = false
)

type IAccount interface {
	Init()
	Run()
	Plot()
	Print(filter ...string)
	Export() [][]string
	analyze()
	filter()
	sort()
}
