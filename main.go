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
	version = "0.0.0"

	flagset   = pflag.NewFlagSet(name, pflag.ContinueOnError)
	isHelp    = flagset.BoolP("help", "h", false, "")
	isVersion = flagset.BoolP("version", "", false, "")
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [OPTION]... [FILE]...
Output direct product of lines of each files.

Options:
      --help                display this help text and exit
      --version             display version information and exit
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

func main() {
	flagset.SetOutput(ioutil.Discard)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printErr(err)
		os.Exit(2)
	}
	if *isHelp {
		printUsage()
		os.Exit(0)
	}
	if *isVersion {
		printVersion()
		os.Exit(0)
	}

	var rs []io.Reader
	if flagset.NArg() == 0 {
		rs = append(rs, os.Stdin)
	} else {
		for _, arg := range flagset.Args() {
			f, err := os.Open(arg)
			if err != nil {
				printErr(err)
				os.Exit(1)
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
			os.Exit(1)
		}
		aa = append(aa, a)
	}

	ss := make([]string, len(aa))
	for indexes := range Product(aa) {
		for i, index := range indexes {
			ss[i] = aa[i][index]
		}
		fmt.Println(strings.Join(ss, "\t"))
	}
}
