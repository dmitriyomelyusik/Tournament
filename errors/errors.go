// Package errors provides usage of custom errors with custom structure
package errors

import (
	"fmt"
)

// Error is the custom made error for convinient handling errors
type Error struct {
	Code    ErrCode     `json:"code"`
	Message string      `json:"message"`
	Info    interface{} `json:"info"`
}

// ErrCode is the code, that used to concritizate error info
type ErrCode string

// Here are all of usable errCodes. Do not delete them!
const (
	RollbackError             ErrCode = "rollbackError"
	CriticalError             ErrCode = "criticalError"
	NotFoundError             ErrCode = "notFoundError"
	DuplicatedIDError         ErrCode = "duplicatedIDError"
	NegativePointsNumberError ErrCode = "negativePointsNumberError"
	NegativeDepositError      ErrCode = "negativeDepositError"
	NoneParticipantsError     ErrCode = "noneParticipantsError"
	ClosedTournamentError     ErrCode = "closedTournamentError"
	UnexpectedError           ErrCode = "unexpectedError"
	JSONError                 ErrCode = "jsonError"
	DatabaseOpenError         ErrCode = "databaseOpenError"
	DatabaseCreatingError     ErrCode = "databaseCreatingError"
	DatabasePingError         ErrCode = "databasePingError"
	TransactionError          ErrCode = "transactionError"
	NotNumberError            ErrCode = "notNumberError"
	ConnectionError           ErrCode = "connectionError"
)

func (e Error) Error() string {
	return fmt.Sprintf("Error: %v\nMessage: %v\nInfo: %v", e.Code, e.Message, e.Info)
}

// Join joins array of errors into one custom error
func Join(errs ...error) Error {
	if len(errs) == 0 {
		return Error{}
	}
	myErr := Transform(errs[0])
	for i := 1; i < len(errs); i++ {
		if errs[i] != nil {
			err := Transform(errs[i])
			myErr.Message += "\n\t" + err.Message
		}
	}
	return myErr
}

// Transform transforms error to Error
func Transform(err error) Error {
	myErr, ok := err.(Error)
	if !ok {
		myErr = Error{
			Code:    UnexpectedError,
			Message: err.Error(),
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
