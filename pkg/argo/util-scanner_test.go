package argo_test

import (
	"fmt"

	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func ExampleDelimitedSliceScanner_commaSeparatedValues() {
	scanner := argo.DelimitedSliceScanner("goodbye,cruel,world", ",")
	values := make([]string, 0, 3)

	for v := range scanner {
		values = append(values, v)
	}

	fmt.Println(values)
	// Output: [goodbye cruel world]
}
