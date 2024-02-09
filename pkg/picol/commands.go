package picol

import (
	errors "errors"
	fmt "fmt"
	strconv "strconv"
	strings "strings"
)

/*
 * incorrectArgCountError returns an error message indicating the incorrect
 * number of arguments provided for a given function. It takes an interpreter
 * instance 'i', the function name 'name', and a slice of argument values 'argv'.
 */
func incorrectArgCountError(i *Interpreter, name string, argv []string) error {
	return fmt.Errorf("wrong number of args for %s %s", name, argv)
}

/*
 * NeedleInHaystack returns true if the string is in a slice
 */
func NeedleInHaystack(needle string, haystack []string) bool {
	for _, haystackMember := range haystack {
		if haystackMember == needle {
			return true
		}
	}
	return false
}

/*
 * CommandMath is the math command for TCL
 */
func CommandMath(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", incorrectArgCountError(i, argv[0], argv)
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
	default:
		return "0", errors.New("invalid operator " + argv[0])
	}
	return fmt.Sprintf("%d", c), nil
}

/*
 * CommandSet is the set command for TCL
 */
func CommandSet(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	i.SetVariable(argv[1], argv[2])
	return argv[2], nil
}

/*
 * CommandUnset is the unset command for TCL
 */
func CommandUnset(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	i.UnsetVariable(argv[1])
	return "", nil
}

/*
 * CommandIf is the if command for TCL
 */
func CommandIf(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 && len(argv) != 5 {
		return "", incorrectArgCountError(i, argv[0], argv)
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

/*
 * CommandWhile is the while command for TCL
 */
func CommandWhile(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 3 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}

	for {
		result, err := i.Eval(argv[1])
		if err != nil {
			return "", err
		}
		if r, _ := strconv.Atoi(result); r != 0 {
			result, err := i.Eval(argv[2])
			switch err {
			case ErrContinue, nil:
				/*
				 * pass
				 */
			case ErrBreak:
				return result, nil
			default:
				return result, err
			}
		} else {
			return result, nil
		}
	}
}

/*
 * CommandRetCodes is a function to get the return codes for TCL
 */
func CommandRetCodes(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	switch argv[0] {
	case "break":
		return "", ErrBreak
	case "continue":
		return "", ErrContinue
	}
	return "", nil
}

/*
 * CommandCallProc is a function to call proc commands for TCL
 */
func CommandCallProc(i *Interpreter, argv []string, pd interface{}) (string, error) {
	var x []string

	if pd, ok := pd.([]string); ok {
		x = pd
	} else {
		return "", nil
	}

	i.callframe = &CallFrame{vars: make(map[string]Variable), parent: i.callframe}
	defer func() { i.callframe = i.callframe.parent }() // remove the called proc callframe

	arity := 0
	for _, arg := range strings.Split(x[0], " ") {
		if len(arg) == 0 {
			continue
		}
		arity++
		i.SetVariable(arg, argv[arity])
	}

	if arity != len(argv)-1 {
		return "", fmt.Errorf("proc '%s' called with wrong arg num", argv[0])
	}

	body := x[1]
	result, err := i.Eval(body)
	if err == ErrReturn {
		err = nil
	}
	return result, err
}

/*
 * CommandProc is a function to register proc commands for TCL
 */
func CommandProc(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 4 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	return "", i.RegisterCommand(argv[1], CommandCallProc, []string{argv[2], argv[3]})
}

/*
 * CommandReturn is a function to register return codes for commands for TCL
 */
func CommandReturn(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	var r string
	if len(argv) == 2 {
		r = argv[1]
	}
	return r, ErrReturn
}

/*
 * CommandError is a function to return error codes for commands for TCL
 */
func CommandError(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 1 && len(argv) != 2 {
		return "", incorrectArgCountError(i, argv[0], argv)
	}
	return "", fmt.Errorf(argv[1])
}

/*
 * CommandPuts is a function to print strings for TCL
 */
func CommandPuts(i *Interpreter, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("wrong number of args for %s %s", argv[0], argv)
	}
	fmt.Println(argv[1])
	return "", nil
}

/*
 * RegisterCoreCommands is a callable to register TCL commands.
 */
func (i *Interpreter) RegisterCoreCommands() error {
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

	return nil
}
