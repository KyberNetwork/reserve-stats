package storage

import (
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/go-pg/pg"
)

func TestNewDB(t *testing.T) {
	type args struct {
		user     string
		password string
		database string
	}
	tests := []struct {
		name string
		args args
		want *UserDB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDB(tt.args.user, tt.args.password, tt.args.database); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDB_StoreUserInfo(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserDB{
				db: tt.fields.db,
			}
			udb.StoreUserInfo()
		})
	}
}

func TestUserDB_GetUserInfo(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.UserUpdate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserDB{
				db: tt.fields.db,
			}
			if got := udb.GetUserInfo(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserDB.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
