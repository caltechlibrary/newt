package main

/*
Newtagain separtes the SQL generation from the newt web service. This is intended to keep things from being confusing about what the `newt` command does.
*/
import (
	"fmt"
	"os"
	"path"
)

func main() {
	appName := path.Base(os.Args[0])
	fmt.Printf("%s will provide the SQL generation previously supplied in the newt prototype\n", appName)
}
