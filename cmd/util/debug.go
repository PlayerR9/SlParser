package util

import "fmt"

func DebugPrint(title string, f func() []string) {
	fmt.Println(title)

	lines := f()

	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Println()
}
