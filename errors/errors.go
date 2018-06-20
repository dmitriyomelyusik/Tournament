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
	PlayerNotFoundError       ErrCode = "playerNotFoundError"
	TournamentNotFoundError   ErrCode = "tournamentNotFoundError"
	DuplicatedIDError         ErrCode = "duplicatedIDError"
	NegativePointsNumberError ErrCode = "negativePointsNumberError"
	NegativeDepositError      ErrCode = "negativeDepositError"
	NoneParticipantsError     ErrCode = "noneParticipantsError"
	ClosedTournamentError     ErrCode = "closedTournamentError"

	UnexpectedError   ErrCode = "unexpectedError"
	JSONError         ErrCode = "jsonError"
	DatabaseOpenError ErrCode = "databaseOpenError"
	TransactionError  ErrCode = "transactionError"
)

func (e Error) Error() string {
	return fmt.Sprintf("Error: %v\nMessage: %v\nInfo: %v", e.Code, e.Message, e.Info)
}
