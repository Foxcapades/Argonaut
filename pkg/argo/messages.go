package argo

import "fmt"

var (
	MultiErrorHeaderLine = func(uniqueErrors, totalErrors int) string {
		return fmt.Sprintf("Encountered %d errors (%d unique):", totalErrors, uniqueErrors)
	}
)
