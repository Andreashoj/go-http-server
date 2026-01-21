package router

import (
	"reflect"
	"testing"
)

func Test_router_FindMatchingRoute(t *testing.T) {
	type fields struct {
		routes []route
	}
	type args struct {
		request HTTPRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				routes: tt.fields.routes,
			}
			if got := r.FindMatchingRoute(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMatchingRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}
