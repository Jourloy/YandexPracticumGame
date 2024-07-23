package errs

import "errors"

type Err struct {
	Error error
	Code  int
}

type errs struct {
	// General

	NotFound      Err
	DatabaseError Err
	AlreadyExist  Err
	ParseError    Err

	// Account

	Unathorized      Err
	InvalidAccount   Err
	NotEnoughBalance Err

	// Building

	PlaceUsed     Err
	TownhallExist Err
}

var (
	NotEnoughBalance = Err{Error: errors.New(`not enough balance`), Code: 400}
)

var Errors = errs{
	// General

	NotFound:      Err{Error: errors.New(`not found`), Code: 404},
	DatabaseError: Err{Error: errors.New(`database error`), Code: 400},
	AlreadyExist:  Err{Error: errors.New(`already exist`), Code: 400},
	ParseError:    Err{Error: errors.New(`parse error`), Code: 400},

	// Account

	Unathorized:    Err{Error: errors.New(`api key is required or account is not exist`), Code: 403},
	InvalidAccount: Err{Error: errors.New(`invalid account`), Code: 400},

	// Building

	PlaceUsed:     Err{Error: errors.New(`place used`), Code: 400},
	TownhallExist: Err{Error: errors.New(`townhall already exist`), Code: 400},
}
