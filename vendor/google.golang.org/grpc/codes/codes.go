



package codes 

import (
	"fmt"
)


type Code uint32

const (
	
	OK Code = 0

	
	Canceled Code = 1

	
	
	
	
	
	Unknown Code = 2

	
	
	
	
	InvalidArgument Code = 3

	
	
	
	
	
	DeadlineExceeded Code = 4

	
	
	NotFound Code = 5

	
	
	AlreadyExists Code = 6

	
	
	
	
	
	
	PermissionDenied Code = 7

	
	
	ResourceExhausted Code = 8

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	FailedPrecondition Code = 9

	
	
	
	
	
	
	Aborted Code = 10

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	OutOfRange Code = 11

	
	
	Unimplemented Code = 12

	
	
	
	Internal Code = 13

	
	
	
	
	
	
	Unavailable Code = 14

	
	DataLoss Code = 15

	
	
	Unauthenticated Code = 16
)

var strToCode = map[string]Code{
	`"OK"`: OK,
	`"CANCELLED"`: Canceled,
	`"UNKNOWN"`:             Unknown,
	`"INVALID_ARGUMENT"`:    InvalidArgument,
	`"DEADLINE_EXCEEDED"`:   DeadlineExceeded,
	`"NOT_FOUND"`:           NotFound,
	`"ALREADY_EXISTS"`:      AlreadyExists,
	`"PERMISSION_DENIED"`:   PermissionDenied,
	`"RESOURCE_EXHAUSTED"`:  ResourceExhausted,
	`"FAILED_PRECONDITION"`: FailedPrecondition,
	`"ABORTED"`:             Aborted,
	`"OUT_OF_RANGE"`:        OutOfRange,
	`"UNIMPLEMENTED"`:       Unimplemented,
	`"INTERNAL"`:            Internal,
	`"UNAVAILABLE"`:         Unavailable,
	`"DATA_LOSS"`:           DataLoss,
	`"UNAUTHENTICATED"`:     Unauthenticated,
}


func (c *Code) UnmarshalJSON(b []byte) error {
	
	
	
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}
	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}
