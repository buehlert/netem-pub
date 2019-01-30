package plus

type PlusParser interface {
	Parse(filename string, filename2 string) (*PlusData, error)
}
