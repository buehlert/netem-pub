package plus

import (
	"fmt"
)

type PlusData struct {
	XData int64
	YData int64
	// SpinServerDelay   float64
	// SpinClientDelay   float64
	// PsnPseServerDelay float64
	// PsnPseClientDelay float64
}

func Parse(text string) (*PlusData, error) {
	var xValue, yValue int64
	fmt.Sscanf(text, "%d,%d",
		&xValue, &yValue)

	return &PlusData{
		XData: xValue,
		YData: yValue,
	}, nil
}
