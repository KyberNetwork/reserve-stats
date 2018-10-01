package stats

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/users/cmc"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

func TestUserStats_GetTxCapByAddress(t *testing.T) {
	type fields struct {
		cmcEthUSDRate *cmc.EthUSDRate
		userStorage   *storage.UserDB
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := UserStats{
				cmcEthUSDRate: tt.fields.cmcEthUSDRate,
				userStorage:   tt.fields.userStorage,
			}
			got, got1, err := us.GetTxCapByAddress(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStats.GetTxCapByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStats.GetTxCapByAddress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserStats.GetTxCapByAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserStats_StoreUserInfo(t *testing.T) {
	type fields struct {
		cmcEthUSDRate *cmc.EthUSDRate
		userStorage   *storage.UserDB
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
			us := UserStats{
				cmcEthUSDRate: tt.fields.cmcEthUSDRate,
				userStorage:   tt.fields.userStorage,
			}
			if err := us.StoreUserInfo(tt.args.email, tt.args.addresses); (err != nil) != tt.wantErr {
				t.Errorf("UserStats.StoreUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUserStats(t *testing.T) {
	type args struct {
		cmc     *cmc.EthUSDRate
		storage *storage.UserDB
	}
	tests := []struct {
		name string
		args args
		want *UserStats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserStats(tt.args.cmc, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
