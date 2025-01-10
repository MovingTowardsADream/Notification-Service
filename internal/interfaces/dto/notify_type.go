package dto

type NotifyType int32

const (
	NotifyTypeModerate    NotifyType = 0
	NotifyTypeSignificant NotifyType = 1
	NotifyTypeAlert       NotifyType = 2
)
