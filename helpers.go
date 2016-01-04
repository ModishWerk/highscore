package main

import (
	"fmt"
	"strconv"
)

func number_to_rank(n int) string {
	switch {
	case n <= 0:
		return ""
	case n == 1:
		return "1st"
	case n == 2:
		return "2nd"
	case n == 3:
		return "3rd"
	default:
		return fmt.Sprintf("%sth", strconv.Itoa(n))
	}
}
