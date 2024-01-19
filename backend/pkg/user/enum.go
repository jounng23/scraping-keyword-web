package user

type Status int32

const (
	Unknown Status = iota
	New
	Approved
	Rejected
)

type Type int32

const (
	TypeUnknown = iota
	TypeReferralPartner
	TypeReferralBuyer
)
