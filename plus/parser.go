package plus

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type PlusData struct {
	PsnPse  float64
	Spin    float64
	Valid   int64
	Invalid int64
	ValidTs int64
}

func Parse(filename string, filename2 string) (*PlusData, error) {
	g, err2 := os.Open(filename2)
	if err2 != nil {
		return nil, err2
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

	toReturn.Valid = currentValid
	toReturn.Invalid = currentInvalid
	toReturn.ValidTs = currentValidTs

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr = csv.NewReader(f)

	var currentSpin, currentPsnPseTemp float64
	currentPsnPse := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
	var newPsnPse, newSpin bool

	for {
		row, err := csvr.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			return &toReturn, err
		}

		if len(row) != 12 {
			continue
		}

		currentPsnPseTemp, err = strconv.ParseFloat(row[10], 64)
		if err != nil {
			return &toReturn, err
		}
		currentSpin, err = strconv.ParseFloat(row[8], 64)
		if err != nil {
			return &toReturn, err
		}
		newPsnPse, err = strconv.ParseBool(row[11])
		if err != nil {
			return &toReturn, err
		}
		newSpin, err = strconv.ParseBool(row[9])
		if err != nil {
			return &toReturn, err
		}

		if newPsnPse {
			if currentPsnPseTemp > 0.300 {
				fOut, err := os.OpenFile("/root/share/debug.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err == nil {

					defer fOut.Close()
					_, _ = fOut.WriteString("PsnPse: ")
					_, _ = fOut.WriteString(row[0])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[1])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[2])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[3])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[4])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[5])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[6])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[7])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[8])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[9])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[10])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[11])
					_, _ = fOut.WriteString("\n")
				}
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
				fOut, err := os.OpenFile("/root/share/debug.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err == nil {

					defer fOut.Close()
					_, _ = fOut.WriteString("Spin: ")
					_, _ = fOut.WriteString(row[0])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[1])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[2])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[3])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[4])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[5])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[6])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[7])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[8])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[9])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[10])
					_, _ = fOut.WriteString(", ")
					_, _ = fOut.WriteString(row[11])
					_, _ = fOut.WriteString("\n")
				}
				continue
			}
			toReturn.Spin = currentSpin
		}

	}

	return &toReturn, nil
}
