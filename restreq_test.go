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
		k string
		v any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]any
	}{
		{
			name: "string test 1",
			args: args{k: "nick", v: any("test")},
			want: map[string]any{"nick": "test"},
		},
		{
			name: "string test 2",
			args: args{k: "nick", v: any("")},
			want: map[string]any{},
		},
		{
			name: "bool test 1",
			args: args{k: "bool", v: true},
			want: map[string]any{"bool": true},
		},
		{
			name: "bool test 2",
			args: args{k: "bool", v: false},
			want: map[string]any{"bool": false},
		},
		{
			name: "float64 test 1",
			args: args{k: "float64", v: 2.34},
			want: map[string]any{"float64": 2.34},
		},
		{
			name: "int32 test 1",
			args: args{k: "int32", v: int64(76)},
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
			r.AddJSONKeyValue(tt.args.k, tt.args.v)

			if !reflect.DeepEqual(r.json, tt.want) {
				t.Errorf("req.JSON() = %v, want %v", r.json, tt.want)
			}
		})
	}
}
