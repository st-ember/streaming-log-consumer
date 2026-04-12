package logretry

type LogRetryStatus string

const (
	StatusSucceeded LogRetryStatus = "succeeded"
	StatusPending   LogRetryStatus = "pending"
	StatusDiscarded LogRetryStatus = "discarded"
)

func (ls LogRetryStatus) IsValidForNew() bool {
	switch ls {
	case StatusPending:
		return true
	case StatusDiscarded:
		return true
	default:
		return false
	}
}
