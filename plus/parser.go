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
	// newLines := 0
	csvr := csv.NewReader(f)

	// var timestamp float64
	var currentPsnPse, currentSpin float64
	var newPsnPse, newSpin bool
	var toReturn PlusData

	for {
		row, err := csvr.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			n++
			return &toReturn, n + nLine, err
		}

		// if n < nLine {
		// 	n++
		// 	continue
		// }

		// timestamp, _ = strconv.ParseFloat(row[0], 64)
		currentPsnPse, err = strconv.ParseFloat(row[6], 64)
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}
		currentSpin, err = strconv.ParseFloat(row[8], 64)
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}
		newPsnPse, err = strconv.ParseBool(row[7])
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}
		newSpin, err = strconv.ParseBool(row[9])
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}

		// fmt.Println(newSpin, row)

		if newPsnPse {
			toReturn.PsnPse = currentPsnPse
		}

		if newSpin {
			toReturn.Spin = currentSpin
		}

		toReturn.NSpin = 2
		n++
		// newLines++
		// if newLines > 3 {
		// 	break
		// }
	}

	// fOut, err := os.OpenFile("/root/share/test_output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// if err == nil {

	// 	defer fOut.Close()

	// 	_, _ = fOut.WriteString(strconv.Itoa(n))
	// 	_, _ = fOut.WriteString("\n")
	// 	_, _ = fOut.WriteString(fmt.Sprintf("%f", toReturn.PsnPse))
	// 	_, _ = fOut.WriteString("\n")
	// }

	return &toReturn, n + nLine, nil
}
