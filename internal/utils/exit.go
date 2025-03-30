package utils

import "os"

var Exit = func(code int) {
	os.Exit(code)
}
