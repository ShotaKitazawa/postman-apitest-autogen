package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilters_Filter(t *testing.T) {
	type args struct {
		testname string
		b        []byte
	}
	tests := []struct {
		name    string
		f       Filters
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test01",
			f: Filters{{
				Testname:          "test1",
				IgnoreJSONPointer: "hoge.fuga",
			}},
			args: args{
				testname: "test1",
				b:        []byte(`{"hoge":{"fuga":"ooo","piyo":"xxx"}}`),
			},
			want:    []byte(`{"hoge":{"piyo":"xxx"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Filter(tt.args.testname, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filters.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("Filters.Filter() return unexpected values:\n%v", diff)
			}
		})
	}
}
