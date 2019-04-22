package gconf

import (
	"flag"
	"fmt"
	"os"
)

func ExamplePrintFlagUsage() {
	fset := flag.NewFlagSet("", flag.ContinueOnError)
	fset.SetOutput(os.Stdout)
	fset.Usage = func() {
		if fset.Name() == "" {
			fmt.Fprintf(fset.Output(), "Usage:\n")
		} else {
			fmt.Fprintf(fset.Output(), "Usage of %s:\n", fset.Name())
		}
		PrintFlagUsage(fset, false)
	}

	fset.String("v", "", "test short string name")
	fset.String("long", "", "test long string name")
	fset.Int("i", 123, "test short int name")
	fset.Int("int", 0, "test long int name")

	fset.Parse([]string{"-h"})

	// Output:
	// Usage:
	//   -i int
	//     	test short int name (default 123)
	//   --int int
	//     	test long int name (default 0)
	//   --long string
	//     	test long string name (default "")
	//   -v string
	//     	test short string name (default "")
}
