package apex

import "errors"

// UnmarshalTo acts as interface to exchange API response
type UnmarshalTo interface {
	GetError() error
}

// Error defines all error information for each request
type Error struct {
	ReturnCode int64  `json:"code"`
	ReturnMsg  string `json:"msg"`
}

// GetError checks and returns an error if it is supplied.
func (e Error) GetError() error {
	if e.ReturnCode != 0 && e.ReturnMsg != "" {
		return errors.New(e.ReturnMsg)
	}
	return nil
}
