package restreq

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func Test_req_JSON(t *testing.T) {
	type fields struct {
		ctx     context.Context
		timeout time.Duration
		url     string
		json    map[string]any
	}
	type args struct {
		i string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]any
	}{
		{
			name: "string test 1",
			args: args{i: "nick=test"},
			want: map[string]any{"nick": "test"},
		},
		{
			name: "string test 2",
			args: args{i: "nick=test=test"},
			want: map[string]any{"nick": "test=test"},
		},
		{
			name: "string test 3",
			args: args{i: "nick="},
			want: map[string]any{},
		},
		{
			name: "string test 4",
			args: args{i: "nick"},
			want: map[string]any{},
		},
		{
			name: "string test 5",
			args: args{i: "=nick=test"},
			want: map[string]any{},
		},
		{
			name: "bool test 1",
			args: args{i: "bool:=true"},
			want: map[string]any{"bool": true},
		},
		{
			name: "bool test 2",
			args: args{i: "bool:=false"},
			want: map[string]any{"bool": false},
		},
		{
			name: "bool test 3",
			args: args{i: "bool:=xxx"},
			want: map[string]any{},
		},
		{
			name: "bool test 4",
			args: args{i: "bool:="},
			want: map[string]any{},
		},
		{
			name: "bool test 5",
			args: args{i: ":=true"},
			want: map[string]any{},
		},
		{
			name: "float64 test 1",
			args: args{i: "float64:=2.34"},
			want: map[string]any{"float64": 2.34},
		},
		{
			name: "int32 test 1",
			args: args{i: "int32:=76"},
			want: map[string]any{"int32": int64(76)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &req{
				ctx:     tt.fields.ctx,
				timeout: tt.fields.timeout,
				url:     tt.fields.url,
				json:    tt.fields.json,
			}
			r.json = make(map[string]any)
			r.SetJSONKeyValue(tt.args.i)

			if !reflect.DeepEqual(r.json, tt.want) {
				t.Errorf("req.JSON() = %v, want %v", r.json, tt.want)
			}
		})
	}
}
