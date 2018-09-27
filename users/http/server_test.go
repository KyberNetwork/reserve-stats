package http

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServer_GetUsers(t *testing.T) {
	type fields struct {
		r    *gin.Engine
		host string
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				r:    tt.fields.r,
				host: tt.fields.host,
			}
			s.GetUsers(tt.args.c)
		})
	}
}

func TestServer_GetUserInfo(t *testing.T) {
	type fields struct {
		r    *gin.Engine
		host string
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				r:    tt.fields.r,
				host: tt.fields.host,
			}
			s.GetUserInfo(tt.args.c)
		})
	}
}

func TestServer_UpdateUserInfo(t *testing.T) {
	type fields struct {
		r    *gin.Engine
		host string
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				r:    tt.fields.r,
				host: tt.fields.host,
			}
			s.UpdateUserInfo(tt.args.c)
		})
	}
}

func TestServer_register(t *testing.T) {
	type fields struct {
		r    *gin.Engine
		host string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				r:    tt.fields.r,
				host: tt.fields.host,
			}
			s.register()
		})
	}
}

func TestServer_Run(t *testing.T) {
	type fields struct {
		r    *gin.Engine
		host string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				r:    tt.fields.r,
				host: tt.fields.host,
			}
			s.Run()
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
