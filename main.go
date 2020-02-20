package main

import (
	"bufio"
	"fileProcessor/core"
	"fmt"
	"os"
	"strings"
)

const extensionsArgPrefix = "-extensions="
const outputFileNamePrefix = "-output="
const pathPrefix = "-path="

func main() {
	sc := bufio.NewScanner(os.Stdin)
	defer sc.Scan()

	out := "concat_out.txt"
	var extensions []string
	var path string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, extensionsArgPrefix) {
			extensions = strings.Split(strings.TrimLeft(arg, extensionsArgPrefix), ";")
			fmt.Printf("extension list set: %v\r\n", extensions)
		}
		if strings.HasPrefix(arg, outputFileNamePrefix) {
			val := strings.TrimLeft(arg, outputFileNamePrefix)
			if len(val) > 0 {
				out = val
			} else {
				fmt.Printf("not null output filename must be set by output= arg")
			}
		}
		if strings.HasPrefix(arg, pathPrefix) {
			val := strings.Trim(strings.TrimLeft(arg, pathPrefix), "\"")
			path = val
		}
	}

	if len(path) == 0 {
		fmt.Printf("must be program argument %v defined...", pathPrefix)
		return
	}
	if len(extensions) == 0 || extensions[0] == "" {
		fmt.Printf("must be program argument %v defined...", extensionsArgPrefix)
		return
	}

	core.
		NewRecursiveProcessor().
		WithHandler(core.ExtensionListValidator(extensions), core.FileConcatenator(out)).
		ProcessPath(path)
}
