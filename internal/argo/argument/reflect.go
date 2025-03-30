package argument

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/foxcapades/argonaut/v3/internal/xreflect"
)

// SiftValidators iterates through all the given validators, sifting them into
// two categories: pre-parse validators and post-parse validators.
//
// Validators are all confirmed to be function types that expect and return the
// correct number and types of values.  If any of the given values does not
// match one of the expected validator function forms, an error will be
// returned.
//
// The valid function forms that are allowed are:
//
//	func(string) error
//	// and
//	func(T, string) error
//
// In the above example, the type T must match the given binding type passed
// as the second argument to this function.
func SiftValidators(validators []any, bindVal *Binding) ([]any, []any, error) {
	includePostParse := bindVal.IsUsable()

	var bt reflect.Type
	if includePostParse {
		bt = xreflect.RootType(reflect.TypeOf(bindVal.Raw))
	} else {
		bt = reflect.TypeOf(nil)
	}

	pre := make([]any, 0, 1)
	post := make([]any, 0, 1)

	for i, fn := range validators {
		count, err := ValidateValidator(i, fn, bt, includePostParse)
		if err != nil {
			return nil, nil, err
		}

		switch count {
		case 1:
			pre = append(pre, fn)
		case 2:
			post = append(post, fn)
		default:
			panic(fmt.Errorf("illegal state: expected count to equal 1 or 2 but it was %d", count))
		}
	}

	return pre, post, nil
}

// ValidateValidator tests the given value (validator) to ensure that it is of a
// valid validator function type.
//
// See SiftValidators for more details about what constitutes a valid validator
// function.
//
// This function returns the number of arguments accepted by the given validator
// function.  If an error is encountered, a value of 0 is returned with that
// error.
func ValidateValidator(
	index int,
	validator any,
	bindingType reflect.Type,
	includePostParse bool,
) (uint8, error) {
	rv := reflect.ValueOf(validator)

	// Verify that the given value is a function.
	if rv.Kind() != reflect.Func {
		return 0, fmt.Errorf("given validator #%d was of invalid type '%s'", index+1, rv.Kind())
	}

	rt := rv.Type()

	// Verify that the given function returns a single param which is an error
	if rt.NumOut() != 1 {
		return 0, fmt.Errorf("given validator #%d returns %d values when 1 was expected", index+1, rt.NumOut())
	}
	if !rt.Out(0).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return 0, fmt.Errorf("given validator #%d returns a value that is not assignable to type 'error'", index+1)
	}

	switch rt.NumIn() {
	case 1:
		return 1, ValidateSoloValidator(index, rt)
	case 2:
		if includePostParse {
			return 2, ValidateDoubleValidator(index, rt, bindingType)
		} else {
			return 0, errors.New("cannot use a post-parse validator with no argument binding")
		}
	default:
		return 0, fmt.Errorf("given validator #%d accepts %d arguments when 1 or 2 were expected", index+2, rt.NumIn())
	}
}

// ValidateSoloValidator tests the given function to ensure the singular
// parameter function accepts the expected type.
func ValidateSoloValidator(idx int, vt reflect.Type) error {
	if !vt.In(0).AssignableTo(reflect.TypeOf("")) {
		return fmt.Errorf("given validator #%d accepts a non-string argument when func(string) error was expected", idx+1)
	}
	return nil
}

// ValidateDoubleValidator tests the given function to ensure the double
// parameter function accepts the expected types.
func ValidateDoubleValidator(idx int, vt reflect.Type, bt reflect.Type) error {
	if !vt.In(0).AssignableTo(bt) {
		return fmt.Errorf("given validator #%d's first parameter (type %s) is incompatible with the argument binding type %s", idx+1, vt.In(0), bt)
	}
	if !vt.In(1).AssignableTo(reflect.TypeOf("")) {
		return fmt.Errorf("given validator #%d's second parameter (type %s) is incompatible with type string", idx+1, vt.In(1))
	}
	return nil
}
