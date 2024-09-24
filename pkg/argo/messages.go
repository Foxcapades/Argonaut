package argo

import "fmt"

var (
	ErrMsgMultiErrorHeaderLine = func(uniqueErrors, totalErrors int) string {
		return fmt.Sprintf("Encountered %d errors (%d unique):", totalErrors, uniqueErrors)
	}

	ErrMsgArgumentDefaultAndRequired = "argument was marked as required, but given a default value"
)
