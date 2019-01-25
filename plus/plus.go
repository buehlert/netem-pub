package plus

// type PlusFetcher interface {
// 	Fetch(iface string) (string, error)
// }

type PlusParser interface {
	Parse(filename string, nLine int) (*PlusData, int, error)
}
