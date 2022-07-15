package bot_errors

import "errors"

var (
	ErrFailedGiveRole     = errors.New("failed to give role")
	ErrFailedTakeRole     = errors.New("failed to take role")
	ErrFailedSendResponse = errors.New("failed to respond to player")
	ErrFailedToRecover    = errors.New("failed to recover from error")
	ErrSomewhereElse      = errors.New("error occured somewhere else and bot failed to recover")
)
