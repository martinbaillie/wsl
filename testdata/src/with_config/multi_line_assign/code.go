package pkg

import "fmt"

func AllowMultiLineAssignCuddle() {
	x := SomeFunc(
		"taking multiple values",
		"suddendly not a single line",
	)
	if x != nil { // want "if statements should only be cuddled with assignments"
		fmt.Println("denied")
	}
}
