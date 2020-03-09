package checker

import (
	"errors"
	"fmt"
	"time"
)

// RetrieveStatusMessage is a help func for retrieving status message from map response of the status server
func RetrieveStatusMessage(m map[string]interface{}) (message *StatusMessage, err error) {
	message = &StatusMessage{}
	for key, value := range m {
		switch key {
		case "status":
			if v, ok := value.(float64); ok {
				if v == 500 {
					message.Status = TimeoutStatus
				} else {
					err = errors.New(fmt.Sprintf("Unknown server status, %f", v))
				}
			} else if v, ok := value.(string); ok {
				message.Status = v
			} else {
				err = errors.New(fmt.Sprintf("Unknown type of status value %s", v))
				return
			}
		case "timestamp":
			if v, ok := value.(float64); ok {
				message.Timestamp = time.Unix(int64(v), 0).Format(time.RFC3339)
			}
		case "message":
			message.Message = fmt.Sprint(value)
		}
	}
	return
}
