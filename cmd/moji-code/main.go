package main

import (
	"fmt"
	"os"

	flag "github.com/saihon/flags"
)

const (
	Name    = "moji-code"
	Version = "v1.0"
	Format  = "%7d %U %-11s %s"
)

type Options struct {
	ranges      bool
	decimal     bool
	hexadecimal bool
	verbose     bool
}

var (
	options  Options
	callback Callback
)

func init() {
	flag.CommandLine.Init(Name, flag.ContinueOnError, false)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: %s [options] [arguments]\n\nOptions:\n", Name)
		flag.PrintCustom()
	}

	callback = func(n uint32, e Entity) {
		if !options.verbose {
			e.Detail = ""
		}
		fmt.Printf(Format+"\n", n, n, e.String, e.Detail)
	}

	flag.Bool("version", 'v', false, "Output version information and exit.\n",
		func(_ flag.Getter) error {
			fmt.Fprintf(flag.CommandLine.Output(), "%s: %s\n", flag.CommandLine.Name(), Version)
			return flag.ErrHelp
		})

	flag.BoolVar(&options.ranges, "range", 'r', false,
		"Specify the range as a decimal,\nhexadecimal or character.\n", nil)

	flag.BoolVar(&options.decimal, "decimal", 'd', false,
		"The specified value is set as a decimal.\n", nil)

	flag.BoolVar(&options.hexadecimal, "hexadecimal", 'x', false,
		"The specified value is set as a hexadecimal.\n", nil)

	flag.BoolVar(&options.verbose, "verbose", 'V', false,
		"Output with the details.\n", nil)

	flag.Bool("control", 'c', false, "Output the Control characters in ASCII code.\n", func(_ flag.Getter) error {
		Each(ASCII.Control, callback)
		return flag.ErrHelp
	})
	flag.Bool("number", 'n', false, "Output the numbers.\n", func(_ flag.Getter) error {
		Each(ASCII.Number, callback)
		return flag.ErrHelp
	})
	flag.Bool("symbol", 's', false, "Output the symbols in ASCII code.\n", func(_ flag.Getter) error {
		Each(ASCII.Symbol, callback)
		return flag.ErrHelp
	})
	flag.Bool("upper-case", 'L', false, "Output the Alphabet upper-case.\n", func(_ flag.Getter) error {
		Each(ASCII.Alphabet.Upper, callback)
		return flag.ErrHelp
	})
	flag.Bool("lower-case", 'l', false, "Output the Alphabet upper-case.\n", func(_ flag.Getter) error {
		Each(ASCII.Alphabet.Lower, callback)
		return flag.ErrHelp
	})
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Recover: %v\n", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}

func _main() int {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if flag.IsIgnorableError(err) {
			return 2
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if flag.NArg() == 0 {
		Each(ASCII.All, callback)
		return 0
	}

	var (
		u32slice Uint32Slice
		err      error
	)
	switch {
	case options.decimal:
		u32slice, err = toUint32Slice(flag.Args(), 10)
	case options.hexadecimal:
		u32slice, err = toUint32Slice(flag.Args(), 16)
	default:
		u32slice, err = toUint32Slice(flag.Args(), -1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if options.ranges {
		table, err := u32slice.ToRangeTable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		Each(table, callback)
		return 0
	}

	u32slice.Each(callback)
	return 0
}
