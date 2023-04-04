package picol_test

import (
	"testing"

	"github.com/dragonheim/gagent/pkg/picol"
)

func TestInterpreter(t *testing.T) {
	interp := picol.NewInterpreter()

	// Register a command
	err := interp.RegisterCommand("test", testCommand, nil)
	if err != nil {
		t.Fatalf("Error registering test command: %v", err)
	}

	// Test command execution
	script := "test hello world"
	result, err := interp.Eval(script)
	if err != nil {
		t.Fatalf("Error executing script: %v", err)
	}
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected result '%s', got '%s'", expected, result)
	}

	// Test variable setting
	interp.SetVariable("x", "42")

	// Test variable retrieval
	val, ok := interp.Variable("x")
	if !ok {
		t.Fatalf("Variable 'x' not found")
	}
	expectedVar := "42"
	if val != picol.Variable(expectedVar) {
		t.Errorf("Expected variable value '%s', got '%s'", expectedVar, val)
	}

	// Test variable unsetting
	interp.UnsetVariable("x")
	_, ok = interp.Variable("x")
	if ok {
		t.Fatalf("Variable 'x' should have been unset")
	}
}

// testCommand is a simple custom command for testing
func testCommand(interp *picol.Interpreter, argv []string, privdata interface{}) (string, error) {
	if len(argv) != 3 {
		return "", nil
	}
	return argv[1] + " " + argv[2], nil
}
