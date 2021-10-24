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

func Parse() Args {
	c := flag.String("collection-file", "", "filepath of postman_collection.json")
	f := flag.String("filter-file", "", "filepath of filter.yaml")
	o := flag.String("output-file", "output.json", "filepath of output_collection.json")
	flag.Parse()
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
