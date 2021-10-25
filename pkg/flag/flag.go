package flag

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	CollectionFile string
	FilterFile     string
	OutputFile     string
}

func Parse(name, version string) Args {
	c := flag.String("collection-file", "", "filepath of postman_collection.json")
	f := flag.String("filter-file", "", "filepath of filter.json")
	o := flag.String("output-file", "output.json", "filepath of output_collection.json")
	v := flag.Bool("version", false, "display version")
	flag.Parse()
	if *v {
		fmt.Fprintf(os.Stdout, "%s %s\n", name, version)
		os.Exit(0)
	}
	a := Args{*c, *f, *o}
	if a.CollectionFile == "" || !exists(a.CollectionFile) {
		fmt.Fprintln(os.Stderr, "-collection-file: invalid arguments")
		os.Exit(1)
	}
	return a
}

func exists(filepath string) bool {
	finfo, err := os.Stat(filepath)
	return err == nil && !finfo.IsDir()
}
