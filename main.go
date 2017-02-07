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
	cmdName    = "prod"
	cmdVersion = "0.1.0"
)

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	separator string
	isHelp    bool
	isVersion bool
}

func NewCLI(stdin io.Reader, stdout io.Writer, stderr io.Writer) *CLI {
	return &CLI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *CLI) printUsage() {
	fmt.Fprintf(c.stderr, `
Usage: %s [OPTION]... [FILE]...
Output direct product of lines of each files.

Options:
  -s, --separator=STRING   use STRING to separate columns (default: \t)
      --help               display this help text and exit
      --version            display version information and exit
`[1:], cmdName)
}

func (c *CLI) printVersion() {
	fmt.Fprintf(c.stderr, "%s\n", cmdVersion)
}

func (c *CLI) printErr(err interface{}) {
	fmt.Fprintf(c.stderr, "%s: %s\n", cmdName, err)
}

func (c *CLI) guideToHelp() {
	fmt.Fprintf(c.stderr, "Try '%s --help' for more information.\n", cmdName)
}

func (c *CLI) parseOptions(args []string) (argFiles []string, err error) {
	flagset := pflag.NewFlagSet(cmdName, pflag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)

	flagset.StringVarP(&c.separator, "separator", "s", "\t", "")
	flagset.BoolVarP(&c.isHelp, "help", "h", false, "")
	flagset.BoolVarP(&c.isVersion, "version", "", false, "")

	if err := flagset.Parse(args); err != nil {
		return nil, err
	}
	return flagset.Args(), nil
}

func (c *CLI) newArgfAsList(argFiles []string) (r []io.Reader, err error) {
	switch len(argFiles) {
	case 0:
		return []io.Reader{c.stdin}, nil
	default:
		rs := make([]io.Reader, len(argFiles))
		for i, argFile := range argFiles {
			f, err := os.Open(argFile)
			if err != nil {
				return nil, err
			}
			rs[i] = f
		}
		return rs, nil
	}
}

func (c *CLI) do(rs []io.Reader) error {
	var aa [][]string
	for _, r := range rs {
		var a []string

		bs := bufio.NewScanner(r)
		for bs.Scan() {
			a = append(a, bs.Text())
		}
		if err := bs.Err(); err != nil {
			return err
		}

		aa = append(aa, a)
	}

	bw := bufio.NewWriter(c.stdout)
	columns := make([]string, len(aa))
	for indexes := range Product(aa) {
		for i, row := range indexes {
			columns[i] = aa[i][row]
		}
		if _, err := bw.WriteString(strings.Join(columns, c.separator) + "\n"); err != nil {
			return err
		}
	}
	return bw.Flush()
}

func (c *CLI) Run(args []string) int {
	argFiles, err := c.parseOptions(args)
	if err != nil {
		c.printErr(err)
		c.guideToHelp()
		return 2
	}
	if c.isHelp {
		c.printUsage()
		return 0
	}
	if c.isVersion {
		c.printVersion()
		return 0
	}

	rs, err := c.newArgfAsList(argFiles)
	if err != nil {
		c.printErr(err)
		c.guideToHelp()
		return 2
	}

	if err = c.do(rs); err != nil {
		c.printErr(err)
		return 1
	}
	return 0
}

func main() {
	c := NewCLI(os.Stdin, os.Stdout, os.Stderr)
	e := c.Run(os.Args[1:])
	os.Exit(e)
}
