package main

import (
	"fmt"

	"github.com/aaron-vaz/{{.Inputs.ProjectName.Value}}/internal"
)

func main() {
	fmt.Println(internal.Message)
}
