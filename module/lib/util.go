
package anyform

import (
	"errors"
	"path/filepath"
	"fmt"
	"os"
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

func WriteFile(filePath string, data []byte, executable bool) error {
	err := MkdirAll(filepath.Dir(filePath))
	if err != nil { return err }
	var perm os.FileMode
	if executable {
		perm = 0770
	} else {
		perm = 0660 
	}
	err = os.WriteFile(filePath, data, perm)
	if err != nil {
		return Errorf("writing file '%v': %w", filePath, err)
	}
	return nil
}

func MkdirAll(path string) error {
  err := os.MkdirAll(path, 0770)
	if err != nil {
		return Errorf("mkdir -p '%v': %w", path, err)
	}
	return nil
}

