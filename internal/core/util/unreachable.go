package util

import (
	"errors"
)

var errAssertFailed = errors.New("assert failed")

// UnreachableErr panics if given err is not nil.
// Serves as an "error handler" for errors that must not happen ever.
func UnreachableErr(err error) {
	UnreachableErrWithMessage(err, "")
}

// UnreachableErrWithMessage is the same as [UnreachableErr],
// but you can provide custom panic message (reason).
func UnreachableErrWithMessage(err error, msg string) {
	if err == nil {
		return
	}
	var text = "unreachable"
	if msg != "" {
		text += ", " + msg
	}
	panic(text + ": " + err.Error())
}

// AssertCond panics if given condition is false.
func AssertCond(cond bool) {
	AssertCondMessage(cond, "")
}

// AssertCondMessage is the same as [AssertCond],
// but you can provide a custom panic message (reason).
func AssertCondMessage(cond bool, msg string) {
	UnreachableErrWithMessage(If(cond, nil, errAssertFailed), msg)
}
