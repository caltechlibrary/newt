package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	appName := path.Base(os.Args[0])
	fmt.Printf("%s will be a Pandoc server API compatible mustache template engine\n", appName)
}
