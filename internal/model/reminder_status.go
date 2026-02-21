package model

//go:generate enumer -type=ReminderStatus -trimprefix=ReminderStatus -transform=snake -json -sql -text -output=reminder_status_enumer.go
type ReminderStatus int

const (
	ReminderStatusPending    ReminderStatus = iota // pending
	ReminderStatusProcessing                       // processing
	ReminderStatusSent                             // sent
	ReminderStatusCancelled                        // cancelled
	ReminderStatusFailed                           // failed
)
