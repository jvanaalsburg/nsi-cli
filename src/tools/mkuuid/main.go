package main

import (
	"fmt"
	"os"
	"strings"
)

func toUuid(id string) string {
	uuid := fmt.Sprintf("%032s", id)
	parts := []string{
		uuid[0:8],
		uuid[8:12],
		uuid[12:16],
		uuid[16:20],
		uuid[20:32],
	}
	return strings.Join(parts, "-")
}

func main() {
	if len(os.Args) == 1 {
		println("Usage: mkuuid <id>")
		return
	}

	fmt.Print(toUuid(os.Args[1]))
}
