package parse

import (
	. "github.com/Foxcapades/Argonaut/v0/internal/log"
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"reflect"
)

func (p *Parser) complete() {
	TraceStart("Parser.complete")
	defer TraceEnd(func() []interface{} { return nil })

	// Assign defaults to positional args
	for _, arg := range p.com.Arguments() {
		p.assignDefault(arg)
	}

	// Assign defaults to flag args
	for _, group := range p.com.FlagGroups() {
		for _, flag := range group.Flags() {
			if flag.HasArgument() {
				p.assignDefault(flag.Argument())
			}
		}
	}

	if p.waiting != nil {
		// waiting may have been hit with defaults
		if p.waiting.Argument().RawValue() == "" {
			if p.isBoolArg(p.waiting.Argument()) {
				util.Must(p.com.Unmarshaler().Unmarshal("true", p.popArg().Binding()))
			} else {
				// TODO: make this a real error
				panic("missing required arg")
			}
		}
	}

	if len(p.reqs) > 0 {
		// TODO: make this a real error
		panic("missing required params")
	}
}


func (p *Parser) assignDefault(arg A.Argument) {
	if !arg.HasDefault() || arg.RawValue() != "" {
		return
	}

	defType := arg.DefaultType()

	if defType.Kind() == reflect.Func {
		defFn := reflect.ValueOf(arg.Default())
		switch defType.NumOut() {
		case 1:
			vals := defFn.Call(nil)
			arg.RootBindValue().Set(vals[0])
			arg.SetRawValue(vals[0].Type().String())
			delete(p.reqs, pointerFor(arg))
			return
		case 2:
			vals := defFn.Call(nil)
			if !vals[1].IsNil() {
				panic(vals[1].Interface())
			}
			if util.IsUnmarshaler(arg.RootBindValue().Type()) {
				arg.RootBindValue().Elem().Set(vals[0])
			} else {
				arg.RootBindValue().Set(vals[0])
			}
			arg.SetRawValue(vals[0].Type().String())
			delete(p.reqs, pointerFor(arg))
			return
		}
		// TODO: Make this a real error, function default
		//       provider did not have the correct number of
		//       return values.
		panic("invalid state")
	}

	if defType.Kind() == reflect.String {
		strVal := arg.RootDefaultValue().String()

		if arg.BindingType().Kind() == reflect.String {
			arg.RootBindValue().Set(arg.RootDefaultValue())
			arg.SetRawValue(strVal)
			delete(p.reqs, pointerFor(arg))
			return
		}

		util.Must(p.com.Unmarshaler().Unmarshal(strVal, arg.Binding()))
		arg.SetRawValue(strVal)
		delete(p.reqs, pointerFor(arg))
		return
	}

	arg.RootBindValue().Set(arg.RootDefaultValue())
	arg.SetRawValue(arg.RootDefaultValue().Type().String())
	delete(p.reqs, pointerFor(arg))
}