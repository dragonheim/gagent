package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	picol "git.dragonheim.net/dragonheim/gagent/src/picol"
)

var fname = flag.String("f", "", "file name")

// CommandPuts is a simple version of the TCL puts function.
func CommandPuts(i *picol.Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("wrong number of args for %s %s", argv[0], argv)
	}
	fmt.Println(argv[1])
	return "", nil
}

func main() {
	flag.Parse()
	interp := picol.InitInterp()
	interp.RegisterCoreCommands()
	interp.RegisterCommand("puts", CommandPuts, nil)

	buf, err := ioutil.ReadFile(*fname)
	if err == nil {
		result, err := interp.Eval(string(buf))
		if err != nil {
			fmt.Println("ERRROR", result, err)
		}
	} else {
		for {
			fmt.Print("picol> ")
			scanner := bufio.NewReader(os.Stdin)
			clibuf, _ := scanner.ReadString('\n')
			result, err := interp.Eval(clibuf[:len(clibuf)-1])
			if len(result) != 0 {
				fmt.Println("ERRROR", result, err)
			}
		}
	}
}
