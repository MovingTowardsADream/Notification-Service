package models

type NotifyType int8

const (
	NotifyTypeModerate    NotifyType = 0
	NotifyTypeSignificant NotifyType = 1
	NotifyTypeAlert       NotifyType = 2
)
