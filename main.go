package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ogier/pflag"
)

var (
	name    = "prod"
	version = "0.1.0"

	flagset   = pflag.NewFlagSet(name, pflag.ContinueOnError)
	separator = flagset.StringP("separator", "s", "\t", "")
	isHelp    = flagset.BoolP("help", "h", false, "")
	isVersion = flagset.BoolP("version", "", false, "")
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [OPTION]... [FILE]...
Output direct product of lines of each files.

Options:
  -s, --separator=STRING   use STRING to separate columns (default: \t)
      --help               display this help text and exit
      --version            display version information and exit
`[1:], name)
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "%s\n", version)
}

func printErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
}

func do(w io.Writer, aa [][]string, separator string) error {
	ss := make([]string, len(aa))
	for indexes := range Product(aa) {
		for i, row := range indexes {
			ss[i] = aa[i][row]
		}
		fmt.Fprintln(w, strings.Join(ss, separator))
	}
	return nil
}

func _main() int {
	flagset.SetOutput(ioutil.Discard)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	if *isHelp {
		printUsage()
		return 0
	}
	if *isVersion {
		printVersion()
		return 0
	}

	var rs []io.Reader
	if flagset.NArg() == 0 {
		rs = append(rs, os.Stdin)
	} else {
		for _, arg := range flagset.Args() {
			f, err := os.Open(arg)
			if err != nil {
				printErr(err)
				guideToHelp()
				return 2
			}
			defer f.Close()
			rs = append(rs, f)
		}
	}

	var aa [][]string
	for _, r := range rs {
		a := make([]string, 0, 64)
		b := bufio.NewScanner(r)
		for b.Scan() {
			a = append(a, b.Text())
		}
		if err := b.Err(); err != nil {
			printErr(err)
			return 1
		}
		aa = append(aa, a)
	}

	if err := do(os.Stdout, aa, *separator); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
