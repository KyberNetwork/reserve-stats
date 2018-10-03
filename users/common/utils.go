package common

//NonKycedCap return cap for non kyc user
func NonKycedCap() *UserCap {
	return &UserCap{
		DailyLimit: 15000.0,
		TxLimit:    3000.0,
	}
}

//KycedCap return cap for kyc user
func KycedCap() *UserCap {
	return &UserCap{
		DailyLimit: 200000.0,
		TxLimit:    6000.0,
	}
}
