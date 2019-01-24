package plus

type PlusFetcher interface {
	Fetch(iface string) (string, error)
}

type PlusParser interface {
	Parse(text string) (*PlusData, error)
}
