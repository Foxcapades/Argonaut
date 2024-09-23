package token

import (
	"fmt"

	"github.com/foxcapades/argonaut/internal/util/xstr"
)

type Token struct {
	Type      Type
	Value     string
	Separator uint16
}

func (t Token) String() string {
	return fmt.Sprintf("Token(Type: %s, Value: %s)", t.Type, xstr.Truncate(t.Value, 16))
}
