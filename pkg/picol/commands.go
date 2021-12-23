package picol

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func arityErr(i *Interp, name string, argv []string) error {
	return fmt.Errorf("wrong number of args for %s %s", name, argv)
}

/*
 needleInHaystack returns true if the string is in a slice
*/
func needleInHaystack(needle string, haystack []string) bool {
	for _, haystackMember := range haystack {
		if haystackMember == needle {
			return true
		}
	}
	return false
}

/*
 TestneedleInHaystack tests the return value of needleInHaystack
*/
func TestneedleInHaystack(t *testing.T) {
	var haystack = []string{"a", "b", "c"}
	var needle = "a"
	if !needleInHaystack(needle, haystack) {
		t.Errorf("%s not in %s", needle, haystack)
	}

	needle = "j"
	if needleInHaystack(needle, haystack) {
		t.Errorf("%s in %s", needle, haystack)
	}

	needle = "ab"
	if needleInHaystack(needle, haystack) {
		t.Errorf("%s in %s", needle, haystack)
	}
}

// CommandMath is the math command for TCL
func CommandMath(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", arityErr(i, argv[0], argv)
	}
	a, _ := strconv.Atoi(argv[1])
	b, _ := strconv.Atoi(argv[2])
	var c int
	switch {
	case argv[0] == "+":
		c = a + b
	case argv[0] == "-":
		c = a - b
	case argv[0] == "*":
		c = a * b
	case argv[0] == "/":
		c = a / b
	case argv[0] == ">":
		if a > b {
			c = 1
		}
	case argv[0] == ">=":
		if a >= b {
			c = 1
		}
	case argv[0] == "<":
		if a < b {
			c = 1
		}
	case argv[0] == "<=":
		if a <= b {
			c = 1
		}
	case argv[0] == "==":
		if a == b {
			c = 1
		}
	case argv[0] == "!=":
		if a != b {
			c = 1
		}
	default: // FIXME I hate warnings
		c = 0
	}
	return fmt.Sprintf("%d", c), nil
}

// CommandSet is the set command for TCL
func CommandSet(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", arityErr(i, argv[0], argv)
	}
	i.SetVar(argv[1], argv[2])
	return argv[2], nil
}

// CommandUnset is the unset command for TCL
func CommandUnset(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", arityErr(i, argv[0], argv)
	}
	i.UnsetVar(argv[1])
	return "", nil
}

// CommandIf is the if command for TCL
func CommandIf(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 && len(argv) != 5 {
		return "", arityErr(i, argv[0], argv)
	}

	result, err := i.Eval(argv[1])
	if err != nil {
		return "", err
	}

	if r, _ := strconv.Atoi(result); r != 0 {
		return i.Eval(argv[2])
	} else if len(argv) == 5 {
		return i.Eval(argv[4])
	}

	return result, nil
}

// CommandWhile is the while command for TCL
func CommandWhile(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", arityErr(i, argv[0], argv)
	}

	for {
		result, err := i.Eval(argv[1])
		if err != nil {
			return "", err
		}
		if r, _ := strconv.Atoi(result); r != 0 {
			result, err := i.Eval(argv[2])
			switch err {
			case errContinue, nil:
				//pass
			case errBreak:
				return result, nil
			default:
				return result, err
			}
		} else {
			return result, nil
		}
	}
}

// CommandRetCodes is a function to get the return codes for TCL
func CommandRetCodes(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 {
		return "", arityErr(i, argv[0], argv)
	}
	switch argv[0] {
	case "break":
		return "", errBreak
	case "continue":
		return "", errContinue
	}
	return "", nil
}

// CommandCallProc is a function to call proc commands for TCL
func CommandCallProc(i *Interp, argv []string, pd interface{}) (string, error) {
	var x []string

	if pd, ok := pd.([]string); ok {
		x = pd
	} else {
		return "", nil
	}

	i.callframe = &CallFrame{vars: make(map[string]Var), parent: i.callframe}
	defer func() { i.callframe = i.callframe.parent }() // remove the called proc callframe

	arity := 0
	for _, arg := range strings.Split(x[0], " ") {
		if len(arg) == 0 {
			continue
		}
		arity++
		i.SetVar(arg, argv[arity])
	}

	if arity != len(argv)-1 {
		return "", fmt.Errorf("proc '%s' called with wrong arg num", argv[0])
	}

	body := x[1]
	result, err := i.Eval(body)
	if err == errReturn {
		err = nil
	}
	return result, err
}

// CommandProc is a function to register proc commands for TCL
func CommandProc(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 4 {
		return "", arityErr(i, argv[0], argv)
	}
	return "", i.RegisterCommand(argv[1], CommandCallProc, []string{argv[2], argv[3]})
}

// CommandReturn is a function to register return codes for commands for TCL
func CommandReturn(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", arityErr(i, argv[0], argv)
	}
	var r string
	if len(argv) == 2 {
		r = argv[1]
	}
	return r, errReturn
}

// CommandError is a function to return error codes for commands for TCL
func CommandError(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", arityErr(i, argv[0], argv)
	}
	return "", fmt.Errorf(argv[1])
}

// CommandPuts is a function to print strings for TCL
func CommandPuts(i *Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("wrong number of args for %s %s", argv[0], argv)
	}
	fmt.Println(argv[1])
	return "", nil
}

// RegisterCoreCommands is a callable to register TCL commands.
func (i *Interp) RegisterCoreCommands() {
	name := [...]string{"+", "-", "*", "/", ">", ">=", "<", "<=", "==", "!="}
	for _, n := range name {
		_ = i.RegisterCommand(n, CommandMath, nil)
	}
	_ = i.RegisterCommand("set", CommandSet, nil)
	_ = i.RegisterCommand("unset", CommandUnset, nil)
	_ = i.RegisterCommand("if", CommandIf, nil)
	_ = i.RegisterCommand("while", CommandWhile, nil)
	_ = i.RegisterCommand("break", CommandRetCodes, nil)
	_ = i.RegisterCommand("continue", CommandRetCodes, nil)
	_ = i.RegisterCommand("proc", CommandProc, nil)
	_ = i.RegisterCommand("return", CommandReturn, nil)
	_ = i.RegisterCommand("error", CommandError, nil)
	_ = i.RegisterCommand("puts", CommandPuts, nil)
}
