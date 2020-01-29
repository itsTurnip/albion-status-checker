package checker

type StatusMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (m StatusMessage) String() string {
	return "<Status: " + m.Status + ". Message: " + m.Message + ">"
}
