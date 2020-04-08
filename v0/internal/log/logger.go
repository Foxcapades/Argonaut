package log

import (
	"fmt"
	"strings"
)

const (
	enabled = false

	prefixSkip = "  | "
	prefixCurr = "  |-"
)

var stack []string

type TraceEndFn = func() []interface{}

func TraceStart(name string, args ...interface{}) {
	if !enabled {
		return
	}

	fmt.Println(prefix(), fmt.Sprintf("Start %s <-- %s", name, args))
	stack = append(stack, name)
}

func Trace(args ...interface{}) {
	if !enabled {
		return
	}
	fmt.Println(prefix(), fmt.Sprint(args...))
}

func Tracef(format string, args ...interface{}) {
	if !enabled {
		return
	}
	fmt.Println(prefix(), fmt.Sprintf(format, args...))
}

func TraceEnd(fn TraceEndFn) {
	if !enabled {
		return
	}
	out := fn()
	name := pop()
	if len(out) == 0 {
		fmt.Println(prefix(), fmt.Sprintf("End %s --> void", name))
	} else {
		fmt.Println(prefix(), fmt.Sprintf("End %s -->", name), out)
	}
}

func pop() (out string) {
	out = stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return
}

func prefix() string {
	if len(stack) == 0 {
		return ""
	}

	out := strings.Builder{}

	for i := 1; i < len(stack); i++ {
		out.WriteString(prefixSkip)
	}

	out.WriteString(prefixCurr)
	return out.String()
}
