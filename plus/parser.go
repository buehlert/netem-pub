package plus

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// type PlusData struct {
// 	XData int64
// 	YData int64
// 	// SpinServerDelay   float64
// 	// SpinClientDelay   float64
// 	// PsnPseServerDelay float64
// 	// PsnPseClientDelay float64
// }

type PlusData struct {
	PsnPse float64
	Spin   float64
	NSpin  float64
}

func Parse(filename string, nLine int) (*PlusData, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nLine, err
	}
	defer f.Close()

	n := 0
	newLines := 0
	csvr := csv.NewReader(f)

	// var timestamp float64
	var currentPsnPse, currentSpin float64
	var newPsnPse, newSpin bool
	var toReturn PlusData

	for {
		row, err := csvr.Read()
		if n < nLine {
			n++
			continue
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			n++
			continue
		}

		// timestamp, _ = strconv.ParseFloat(row[0], 64)
		currentPsnPse, _ = strconv.ParseFloat(row[6], 64)
		currentSpin, _ = strconv.ParseFloat(row[8], 64)
		newPsnPse, _ = strconv.ParseBool(row[7])
		newSpin, _ = strconv.ParseBool(row[9])

		// fmt.Println(newSpin, row)

		if newPsnPse {
			toReturn.PsnPse = currentPsnPse
		}

		if newSpin {
			toReturn.Spin = currentSpin
		}

		toReturn.NSpin = 2
		n++
		newLines++
		if newLines > 3 {
			break
		}
	}

	return &toReturn, n, nil
}
