package checker

const OnlineStatus = "online"
const OfflineStatus = "offline"
const TimeoutStatus = "timeout"
const StartingStatus = "starting"

// Status message represents Albion server response for status checks.
type StatusMessage struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// String returns string representation of StatusMessage
func (m *StatusMessage) String() string {
	return "<Status: " + m.Status + ". Message: " + m.Message + ">"
}
