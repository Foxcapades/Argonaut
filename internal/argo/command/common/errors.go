package common

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func AppendError(result argo.ParseResult, err error) (argo.ParseResult, error) {
	result.Error = err
	return result, err
}

func AppendWarning(result *argo.ParseResult, warning string, kind argo.ParseWarningType) {
	result.Warnings = append(result.Warnings, argo.ParseWarning{
		Type:    kind,
		Message: warning,
	})
}

func AppendUnrecognizedLongFlag(result *argo.ParseResult, name string) {
	AppendWarning(result, fmt.Sprintf("unrecognized long flag --%s", name), argo.UnrecognizedFlag)
}

func AppendUnrecognizedShortFlag(result *argo.ParseResult, name byte) {
	AppendWarning(result, fmt.Sprintf("unrecognized short flag -%c", name), argo.UnrecognizedFlag)
}
