package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	substitution  Substitution
	file          *os.File
	lineRange     string
	doubleSpacing bool
	editInPlace   bool
}

type Substitution struct {
	command     string
	pattern     *regexp.Regexp
	replacement string
	flag        string
}

func main() {
	cfg := loadConfig()

	substitute(cfg)

}

func loadConfig() Config {
	var err error
	var cfg Config

	// output a range of lines from the file. specify a range, i.e.
	// for lines 2 to 4 we would use the command: cat -n ccsed -n '2,4p’ filename
	flag.StringVar(&cfg.lineRange, "n", "", "output a range of lines from the file")
	flag.BoolVar(&cfg.doubleSpacing, "G", false, "double spacing a file")
	flag.BoolVar(&cfg.editInPlace, "i", false, "edit in place")

	flag.Parse()
	args := flag.Args()

	subst := ""

	if len(args) < 1 {
		fmt.Println("Please provide substitution and file as args")
		os.Exit(1)
	} else if len(args) < 2 {
		subst = args[0]
		cfg.file = os.Stdin
	} else {
		subst = args[0]
		cfg.file, err = os.Open(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cfg.substitution, err = parseSubstitution(subst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return cfg
}

func parseSubstitution(subst string) (Substitution, error) {
	var res Substitution
	var err error

	substList := strings.Split(subst, "/")
	if len(substList) != 4 {
		err = fmt.Errorf("invalid substitution: %v", subst)
		return res, err
	}

	res.command = substList[0]
	res.pattern, err = regexp.Compile(substList[1])
	if err != nil {
		err = fmt.Errorf("Not valid regex pattern: %v", substList[1])
		return res, err
	}
	res.replacement = substList[2]
	res.flag = substList[3]

	//fmt.Printf("Parsed substitution:\n%v\n%v\n%v\n%v", res.command, res.pattern, res.replacement, res.flag)

	return res, nil

}

func substitute(cfg Config) {
	re := cfg.substitution.pattern
	repl := cfg.substitution.replacement
	scanner := bufio.NewScanner(cfg.file)

	for scanner.Scan() {
		fmt.Println(re.ReplaceAllString(scanner.Text(), repl))
	}
}
