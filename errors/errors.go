// Package errors provides using custom errors with convinient structure
package errors

import (
	"fmt"
)

// Error is the custom made error for convinient handling errors
type Error struct {
	Code    ErrCode     `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}

// ErrCode is the code, that used to concritizate error info
type ErrCode string

// Here are all of usable errCodes. Do not delete them!
const (
	DatabaseOpenError         ErrCode = "databaseOpenError"
	UnexpectedError           ErrCode = "unexpectedError"
	PlayerNotFoundError       ErrCode = "playerNotFoundError"
	NegativePointsNumberError ErrCode = "negativePointsNumberError"
	DuplicatedIDError         ErrCode = "duplicatedIDError"
	JSONError                 ErrCode = "jsonError"
	TransactionError          ErrCode = "transactionError"
	NegativeDepositError      ErrCode = "negativeDepositError"
	TournamentNotFoundError   ErrCode = "tournamentNotFoundError"
)

func (e Error) Error() string {
	return fmt.Sprintf("Error: %v\nMessage: %v\nInfo: %v", e.Code, e.Message, e.Info)
}
