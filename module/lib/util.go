
package anyform

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// Error wrapper ///////////////////////////////////////////////////////////////

// Tags along with another error to keep track of where the stack trace happened
type errorWithStack struct {
	stackTrace string
}

func (ews *errorWithStack) Error() string {
	return fmt.Sprintf("stack trace:\n%v", ews.stackTrace)
}

func (ews *errorWithStack) Is(target error) bool {
	_, ok := target.(*errorWithStack)
	return ok
}

// Current implementation: effectively the same as fmt.Errorf except that it
// adds a stack trace at the end of the original message, but only if there
// isn't already a stack trace included.  This means the following will only
// include 1 stack trace, and it will be based on the location of the inner most
// Errorf() call:
// 	Errorf("foo: %w", fmt.Errorf("bar: %w", Errorf("baz: %w", origErr)))
func Errorf(format string, arg ...any) error {
	origErr := fmt.Errorf(format, arg...)
	if errors.Is(origErr, &errorWithStack{}) {
		return origErr
	} else {
		return fmt.Errorf(
				"%w\n%w",
				origErr,
				&errorWithStack{ stackTrace: string(debug.Stack()) })
	}
}

func WrapError(err error) error {
	return fmt.Errorf(
		"%w\n%w",
		err,
		&errorWithStack{stackTrace: string(debug.Stack())},
	)
}
