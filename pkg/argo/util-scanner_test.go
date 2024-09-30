package argo_test

import (
	"fmt"

	"github.com/foxcapades/argonaut/pkg/argo"
)

func ExampleDelimitedSliceScanner_commaSeparatedValues() {
	scanner := argo.DelimitedSliceScanner("goodbye,cruel,world", ",")
	values := make([]string, 0, 3)

	for scanner.HasNext() {
		values = append(values, scanner.Next())
	}

	fmt.Println(values)

	// Output: [goodbye cruel world]
}
