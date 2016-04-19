package main

import (
	"bufio"
	"fmt"
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

type productor struct {
	items   [][]string
	indexes []int
	ch      chan []int
}

func newProductor(items [][]string, ch chan []int) *productor {
	return &productor{
		items:   items,
		indexes: make([]int, len(items)),
		ch:      ch,
	}
}

func (p *productor) findNext(index_i int) {
	if index_i == len(p.items) {
		indexes := make([]int, len(p.indexes))
		copy(indexes, p.indexes)
		p.ch <- indexes
		return
	}

	for i := 0; i < len(p.items[index_i]); i++ {
		p.indexes[index_i] = i
		p.findNext(index_i + 1)
	}
}

func Product(items [][]string) chan []int {
	ch := make(chan []int, 16)
	go func() {
		p := newProductor(items, ch)
		p.findNext(0)
		close(p.ch)
	}()
	return ch
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

	var a [][]string
	if flagset.NArg() < 1 {
		a = append(a, make([]string, 0, 64))

		b := bufio.NewScanner(os.Stdin)
		for b.Scan() {
			a[0] = append(a[0], b.Text())
		}
		if err := b.Err(); err != nil {
			printErr(err)
			os.Exit(1)
		}
	} else {
		for i, arg := range flagset.Args() {
			a = append(a, make([]string, 0, 64))

			f, err := os.Open(arg)
			if err != nil {
				printErr(err)
				os.Exit(1)
			}
			defer f.Close()

			b := bufio.NewScanner(f)
			for b.Scan() {
				a[i] = append(a[i], b.Text())
			}
			if err = b.Err(); err != nil {
				printErr(err)
				os.Exit(1)
			}
		}
	}

	ss := make([]string, len(a))
	for indexes := range Product(a) {
		for i, index := range indexes {
			ss[i] = a[i][index]
		}
		fmt.Println(strings.Join(ss, "\t"))
	}
}
