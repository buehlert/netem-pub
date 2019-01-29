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
	PsnPse  float64
	Spin    float64
	NSpin   float64
	Valid   int64
	Invalid int64
	ValidTs int64
}

func Parse(filename string, filename2 string, nLine int) (*PlusData, int, error) {
	g, err2 := os.Open(filename2)
	if err2 != nil {
		return nil, nLine, err2
	}
	defer g.Close()

	csvr := csv.NewReader(g)

	var toReturn PlusData
	var currentValid int64
	var currentInvalid int64
	var currentValidTs int64

	for {
		row, err := csvr.Read()
		if err != nil {
			break
		}

		if len(row) != 3 {
			continue
		}

		currentValid, err = strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			break
		}
		currentInvalid, err = strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			break
		}
		currentValidTs, err = strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			break
		}
	}

	fOut, err3 := os.OpenFile("/root/share/test_output_valid.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err3 == nil {

		defer fOut.Close()

		_, _ = fOut.WriteString(strconv.FormatInt(currentValid, 10))
		_, _ = fOut.WriteString("\n")
		_, _ = fOut.WriteString(strconv.FormatInt(currentInvalid, 10))
		_, _ = fOut.WriteString("\n")
	}

	toReturn.Valid = currentValid
	toReturn.Invalid = currentInvalid
	toReturn.ValidTs = currentValidTs

	f, err := os.Open(filename)
	if err != nil {
		return nil, nLine, err
	}
	defer f.Close()

	n := 0
	// newLines := 0
	csvr = csv.NewReader(f)

	// var timestamp float64
	var currentSpin, currentPsnPseTemp float64
	currentPsnPse := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
	var newPsnPse, newSpin bool

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

		if len(row) != 12 {
			n++
			continue
		}

		// timestamp, _ = strconv.ParseFloat(row[0], 64)

		currentPsnPseTemp, err = strconv.ParseFloat(row[10], 64)
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}
		currentSpin, err = strconv.ParseFloat(row[8], 64)
		if err != nil {
			n++
			return &toReturn, n + nLine, err
		}
		newPsnPse, err = strconv.ParseBool(row[11])
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
			if currentPsnPseTemp > 0.300 {
				n++
				continue
			}
			currentPsnPse[4] = currentPsnPse[3]
			currentPsnPse[3] = currentPsnPse[2]
			currentPsnPse[2] = currentPsnPse[1]
			currentPsnPse[1] = currentPsnPse[0]
			currentPsnPse[0] = currentPsnPseTemp

			toReturn.PsnPse = (currentPsnPse[0] + currentPsnPse[1] + currentPsnPse[2] + currentPsnPse[3] + currentPsnPse[4]) / 5
		}

		if newSpin {
			if currentSpin > 0.250 {
				n++
				continue
			}
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
