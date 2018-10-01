package storage

import (
	"testing"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/go-pg/pg"
)

func newConnectionToTestDB() {

}

func TestUserDB_StoreUserInfo(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		email     string
		addresses []common.UserAddress
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserDB{
				db: tt.fields.db,
			}
			if err := udb.StoreUserInfo(tt.args.email, tt.args.addresses); (err != nil) != tt.wantErr {
				t.Errorf("UserDB.StoreUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
