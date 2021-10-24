package model

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestset_ToString(t *testing.T) {
	type fields struct {
		name        string
		status      int
		Evaluations []*Evaluation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				name:   "test1",
				status: 200,
				Evaluations: []*Evaluation{
					{
						Left:  []Field{{"hoge", false}},
						Right: "fuga",
					},
					{
						Left:  []Field{{"a", false}, {"xxx", false}},
						Right: "aaa",
					},
					{
						Left:  []Field{{"b", false}, {"xxx", false}},
						Right: "bbb",
					},
				},
			},
			want: `pm.test("test1", function(){ responseCode.code === 200; var data = pm.response.json(); pm.expect(data.hoge).to.eql("fuga"); pm.expect(data.xxx.a).to.eql("aaa"); pm.expect(data.xxx.b).to.eql("bbb"); })`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Testset{
				Name:        tt.fields.name,
				Status:      tt.fields.status,
				Evaluations: tt.fields.Evaluations,
			}
			if diff := cmp.Diff(tt.want, tr.ToString()); diff != "" {
				t.Errorf("Testset.String() return unexpected values:\n%v", diff)
			}
		})
	}
}

func TestTestset_parseTest(t *testing.T) {
	type fields struct {
		Evaluations []*Evaluation
	}
	type args struct {
		jsonStr []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Evaluation
	}{
		{
			name: "test1",
			args: args{
				jsonStr: []byte(`{"hoge": "fuga", "xxx": {"a": "aaa", "b": "bbb"}}`),
			},
			want: []*Evaluation{
				{
					Left:  []Field{{"hoge", false}},
					Right: "fuga",
				},
				{
					Left:  []Field{{"a", false}, {"xxx", false}},
					Right: "aaa",
				},
				{
					Left:  []Field{{"b", false}, {"xxx", false}},
					Right: "bbb",
				},
			},
		},
		{
			name:   "test2",
			fields: fields{},
			args: args{
				jsonStr: []byte(`[1, 2, 3]`)},
			want: []*Evaluation{
				{
					Left:  []Field{{"0", true}},
					Right: float64(1),
				},
				{
					Left:  []Field{{"1", true}},
					Right: float64(2),
				},
				{
					Left:  []Field{{"2", true}},
					Right: float64(3),
				},
			},
		},
		{
			name:   "test3",
			fields: fields{},
			args: args{
				jsonStr: []byte(`[{"hoge": "fuga"}, {"hoge": "piyo"}]`),
			},
			want: []*Evaluation{
				{
					Left:  []Field{{"hoge", false}, {"0", true}},
					Right: "fuga",
				},
				{
					Left:  []Field{{"hoge", false}, {"1", true}},
					Right: "piyo",
				},
			},
		},
		{
			name:   "test4",
			fields: fields{},
			args: args{
				jsonStr: []byte(`[{"hoge": [1, 2, 3]}, {"hoge": [4, 5, 6]}]`)},
			want: []*Evaluation{
				{
					Left:  []Field{{"0", true}, {"hoge", false}, {"0", true}},
					Right: float64(1),
				},
				{
					Left:  []Field{{"1", true}, {"hoge", false}, {"0", true}},
					Right: float64(2),
				},
				{
					Left:  []Field{{"2", true}, {"hoge", false}, {"0", true}},
					Right: float64(3),
				},
				{
					Left:  []Field{{"0", true}, {"hoge", false}, {"1", true}},
					Right: float64(4),
				},
				{
					Left:  []Field{{"1", true}, {"hoge", false}, {"1", true}},
					Right: float64(5),
				},
				{
					Left:  []Field{{"2", true}, {"hoge", false}, {"1", true}},
					Right: float64(6),
				},
			},
		},
		{
			name: "test5",
			args: args{
				jsonStr: []byte(`[{"hoge": [[1, 2], [3]]}, {"hoge": [[4], [5, 6]]}]`)},
			want: []*Evaluation{
				{
					Left:  []Field{{"0", true}, {"0", true}, {"hoge", false}, {"0", true}},
					Right: float64(1),
				},
				{
					Left:  []Field{{"1", true}, {"0", true}, {"hoge", false}, {"0", true}},
					Right: float64(2),
				},
				{
					Left:  []Field{{"0", true}, {"1", true}, {"hoge", false}, {"0", true}},
					Right: float64(3),
				},
				{
					Left:  []Field{{"0", true}, {"0", true}, {"hoge", false}, {"1", true}},
					Right: float64(4),
				},
				{
					Left:  []Field{{"0", true}, {"1", true}, {"hoge", false}, {"1", true}},
					Right: float64(5),
				},
				{
					Left:  []Field{{"1", true}, {"1", true}, {"hoge", false}, {"1", true}},
					Right: float64(6),
				},
			},
		},
	}
	for _, tt := range tests {
		var data interface{}
		if err := json.Unmarshal(tt.args.jsonStr, &data); err != nil {
			panic(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			tr := Testset{
				Evaluations: tt.fields.Evaluations,
			}
			if diff := cmp.Diff(tt.want, tr.parseTest(data)); diff != "" {
				t.Errorf("Testset.parseTest() return unexpected values:\n%v", diff)
			}
		})
	}
}
