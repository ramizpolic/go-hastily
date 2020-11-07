package global

import (
	"fmt"
	"os"
	"reflect"
)

// HandleError controls how CLI commands react to an exception.
func HandleError(err error) {
	if err != nil {
		errorMsg := fmt.Sprintf("%v", err)
		if errorMsg != "" {
			CLI.Error("%s", errorMsg)
		} else {
			CLI.Error("Something went wrong.")
		}
		os.Exit(1)
	}
}

// HandleErrorMessage controls how CLI commands react to an exception.
func HandleErrorMessage(err string) {
	if err != "" {
		CLI.Error(err)
		os.Exit(1)
	}
}

// IsZero checks if values are uninitialized.
func IsZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// PrintNonZero prints only non-zero elements.
func PrintNonZero(x interface{}) string {
	if !IsZero(x) {
		return fmt.Sprintf("%v", x)
	}
	return ""
}
