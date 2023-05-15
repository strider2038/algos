package testcases

type Case struct {
	Name   string
	Search string
	Want   []string
}

var SimpleCases = []Case{
	{
		Name: "empty search",
		Want: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"},
	},
	{
		Name:   "top node filter, single word",
		Search: "рыболовство",
		Want:   []string{"A"},
	},
	{
		Name:   "deep node filter, single word",
		Search: "картофел",
		Want:   []string{"A", "C", "G"},
	},
	{
		Name:   "deep node filter, single word with stemming",
		Search: "картофель",
		Want:   []string{"A", "C", "G"},
	},
	{
		Name:   "deep node filter, many words",
		Search: "Деятельность прочих общественных и прочих некоммерческих организаций, кроме религиозных и политических организаций, не включенных в другие группировки",
		Want:   []string{"S"},
	},
	{
		Name:   "single char search",
		Search: "а",
		Want:   []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"},
	},
	{
		Name:   "prefix search",
		Search: "кар",
		Want:   []string{"A", "B", "C", "E", "G", "J", "M"},
	},
	{
		Name:   "stop word as prefix",
		Search: "без",
		Want:   []string{"C", "G", "M", "N", "O", "P", "Q"},
	},
	{
		Name:   "stop word",
		Search: "вдруг",
		Want:   []string{},
	},
}

// PrefixCases - для префиксных алгоритмов некоторые тестовые кейсы отличаются от
// простых случаев (поиск идет не по вхождению в любую часть подстроки,
// а по префиксу).
var PrefixCases = adoptForPrefixCases(SimpleCases)

func adoptForPrefixCases(cases []Case) []Case {
	adopted := make([]Case, 0, len(cases))

	for _, c := range cases {
		switch c.Name {
		case "single char search":
			c.Want = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R"}
		case "stop word as prefix":
			c.Want = []string{"C", "G", "M", "N", "O"}
		}
		adopted = append(adopted, c)
	}

	return adopted
}
