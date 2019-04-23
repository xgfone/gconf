package gconf

import (
	"flag"
	"fmt"
	"os"
)

func ExamplePrintFlagUsage() {
	fset := flag.NewFlagSet("app", flag.ContinueOnError)
	fset.Usage = func() {
		/// For Go1.10+
		// if fset.Name() == "" {
		//     fmt.Fprintf(fset.Output(), "Usage:\n")
		// } else {
		//     fmt.Fprintf(fset.Output(), "Usage of %s:\n", fset.Name())
		// }
		// PrintFlagUsage(fset.Output(), fset, false)

		// Here only for test
		fmt.Fprintf(os.Stdout, "Usage: app\n")
		PrintFlagUsage(os.Stdout, fset, false)
	}

	fset.String("v", "", "test short string name")
	fset.String("long", "", "test long string name")
	fset.Int("i", 123, "test short int name")
	fset.Int("int", 0, "test long int name")

	fset.Parse([]string{"-h"})

	// Output:
	// Usage: app
	//   -i int
	//     	test short int name (default 123)
	//   --int int
	//     	test long int name (default 0)
	//   --long string
	//     	test long string name (default "")
	//   -v string
	//     	test short string name (default "")
}
