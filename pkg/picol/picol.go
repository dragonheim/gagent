package picol

import (
	"errors"
	"fmt"
	"strings"
)

/*
 * Error variables
 */
var (
	ErrReturn   = errors.New("RETURN")
	ErrBreak    = errors.New("BREAK")
	ErrContinue = errors.New("CONTINUE")
)

/*
 * Variable type
 */
type Variable string

/*
 * CommandFunc type
 */
type CommandFunc func(interp *Interpreter, argv []string, privdata interface{}) (string, error)

/*
 * Command structure
 */
type Command struct {
	fn       CommandFunc
	privdata interface{}
}

/*
 * CallFrame structure
 */
type CallFrame struct {
	vars   map[string]Variable
	parent *CallFrame
}

/*
 * Interpreter structure
 */
type Interpreter struct {
	level     int
	callframe *CallFrame
	commands  map[string]Command
}

/*
 * NewInterpreter initializes a new Interpreter
 */
func NewInterpreter() *Interpreter {
	return &Interpreter{
		level:     0,
		callframe: &CallFrame{vars: make(map[string]Variable)},
		commands:  make(map[string]Command),
	}
}

/*
 * Variable retrieves a variable's value
 */
func (interp *Interpreter) Variable(name string) (Variable, bool) {
	for frame := interp.callframe; frame != nil; frame = frame.parent {
		v, ok := frame.vars[name]
		if ok {
			return v, ok
		}
	}
	return "", false
}

/*
 * SetVariable sets a variable's value
 */
func (interp *Interpreter) SetVariable(name, val string) {
	interp.callframe.vars[name] = Variable(val)
}

/*
 * UnsetVariable removes a variable
 */
func (interp *Interpreter) UnsetVariable(name string) {
	delete(interp.callframe.vars, name)
}

/*
 * Command retrieves a command
 */
func (interp *Interpreter) Command(name string) *Command {
	v, ok := interp.commands[name]
	if !ok {
		return nil
	}
	return &v
}

/*
 * RegisterCommand registers a new command
 */
func (interp *Interpreter) RegisterCommand(name string, fn CommandFunc, privdata interface{}) error {
	cmd := interp.Command(name)
	if cmd != nil {
		return fmt.Errorf("Command '%s' already defined", name)
	}

	interp.commands[name] = Command{fn, privdata}
	return nil
}

/*
 * Eval evaluates a script
 */
func (interp *Interpreter) Eval(script string) (string, error) {
	parser := InitParser(script)
	var result string
	var err error

	argv := []string{}

	for {
		prevType := parser.Type
		token := parser.GetToken()
		if parser.Type == ParserTokenEOF {
			break
		}

		switch parser.Type {
		case ParserTokenVAR:
			v, ok := interp.Variable(token)
			if !ok {
				return "", fmt.Errorf("no such variable '%s'", token)
			}
			token = string(v)
		case ParserTokenCMD:
			result, err = interp.Eval(token)
			if err != nil {
				return result, err
			}
			token = result
		case ParserTokenESC:
			/*
			 * TODO: escape handling missing!
			 */
		case ParserTokenSEP:
			// prevType = parser.Type
			continue
		}

		if parser.Type == ParserTokenEOL {
			// prevType = parser.Type
			if len(argv) != 0 {
				cmd := interp.Command(argv[0])
				if cmd == nil {
					return "", fmt.Errorf("no such command '%s'", argv[0])
				}
				result, err = cmd.fn(interp, argv, cmd.privdata)
				if err != nil {
					return result, err
				}
			}
			/*
			 * Prepare for the next command
			 */
			argv = []string{}
			continue
		}

		/*
		 * We have a new token, append to the previous or as new arg?
		 */
		if prevType == ParserTokenSEP || prevType == ParserTokenEOL {
			argv = append(argv, token)
		} else { // Interpolation
			argv[len(argv)-1] = strings.Join([]string{argv[len(argv)-1], token}, "")
		}
		// prevType = parser.Type
	}
	return result, nil
}
