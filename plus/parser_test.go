package plus

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		filename string
		nLine    int
	}
	tests := []struct {
		name    string
		args    args
		want    *PlusData
		wantN   int
		wantErr bool
	}{
		{
			"Test 1",
			args{"test_small.csv", 0},
			&PlusData{PsnPse: 0.086256, Spin: 0.03232, NSpin: 2},
			3,
			false,
		},
		{
			"Test 2",
			args{"test_small.csv", 2},
			&PlusData{PsnPse: 0.086256, Spin: 0.03232, NSpin: 2},
			5,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotN, err := Parse(tt.args.filename, tt.args.nLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Parse() nLine = %v, wantN = %v", gotN, tt.wantN)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
