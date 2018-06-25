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
	NotFoundError             ErrCode = "notFoundError"
	DuplicatedIDError         ErrCode = "duplicatedIDError"
	NegativePointsNumberError ErrCode = "negativePointsNumberError"
	NegativeDepositError      ErrCode = "negativeDepositError"
	NoneParticipantsError     ErrCode = "noneParticipantsError"
	ClosedTournamentError     ErrCode = "closedTournamentError"
	UnexpectedError           ErrCode = "unexpectedError"
	JSONError                 ErrCode = "jsonError"
	DatabaseOpenError         ErrCode = "databaseOpenError"
	TransactionError          ErrCode = "transactionError"
	NotNumberError            ErrCode = "notNumberError"
)

func (e Error) Error() string {
	return fmt.Sprintf("Error: %v\nMessage: %v\nInfo: %v", e.Code, e.Message, e.Info)
}

// Join joins array of errors into one custom error
func Join(errs ...error) Error {
	if len(errs) == 0 {
		return Error{}
	}
	myErr, ok := errs[0].(Error)
	if !ok {
		myErr = Error{
			Code:    UnexpectedError,
			Message: errs[0].Error(),
		}
	}
	for i := 1; i < len(errs); i++ {
		if errs[i] != nil {
			myErr.Message += "\n" + errs[i].Error()
		}
	}
	return myErr
}

// SetPrefix sets prefix in err message
func (e Error) SetPrefix(pref string) Error {
	e.Message = pref + e.Message
	return e
}

// SetCode sets ErrCode
func (e Error) SetCode(code ErrCode) Error {
	e.Code = code
	return e
}
