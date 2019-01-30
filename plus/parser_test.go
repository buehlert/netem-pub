package plus

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		filename  string
		filename2 string
	}
	tests := []struct {
		name    string
		args    args
		want    *PlusData
		wantErr bool
	}{
		{
			"Test 1",
			args{"test_small.csv", "test_valid.csv"},
			&PlusData{PsnPse: (0.0423 + 0.04456 + 0.04216 + 0.04435 + 0.0443) / 5, Spin: 0.03232, Valid: 32, Invalid: 1040, ValidTs: 4567},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.filename, tt.args.filename2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
