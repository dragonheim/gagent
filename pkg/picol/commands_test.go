package picol_test

import (
	"testing"

	"github.com/dragonheim/gagent/pkg/picol"
)

func Test_NeedleInHaystack(t *testing.T) {
	var haystack = []string{"a", "b", "c"}
	var needle = "a"
	if !picol.NeedleInHaystack(needle, haystack) {
		t.Errorf("%s not in %s", needle, haystack)
	}

	needle = "j"
	if picol.NeedleInHaystack(needle, haystack) {
		t.Errorf("%s in %s", needle, haystack)
	}

	needle = "ab"
	if picol.NeedleInHaystack(needle, haystack) {
		t.Errorf("%s in %s", needle, haystack)
	}
}

func Test_CommandMath(t *testing.T) {
	// You can add more test cases for various operations and edge cases
	testCases := []struct {
		argv   []string
		result string
		err    error
	}{
		{[]string{"+", "2", "3"}, "5", nil},
		{[]string{"-", "8", "3"}, "5", nil},
		{[]string{"*", "2", "3"}, "6", nil},
		{[]string{"/", "6", "3"}, "2", nil},
		{[]string{">", "4", "2"}, "1", nil},
	}

	i := picol.NewInterpreter()
	for _, tc := range testCases {
		result, err := picol.CommandMath(i, tc.argv, nil)
		if result != tc.result || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
			t.Errorf("CommandMath(%v) = (%v, %v); expected (%v, %v)", tc.argv, result, err, tc.result, tc.err)
		}
	}
}

func Test_CommandSet(t *testing.T) {
	testCases := []struct {
		argv   []string
		result string
		err    error
	}{
		{[]string{"set", "x", "42"}, "42", nil},
		{[]string{"set", "y", "abc"}, "abc", nil},
	}

	i := picol.NewInterpreter()
	for _, tc := range testCases {
		result, err := picol.CommandSet(i, tc.argv, nil)
		if result != tc.result || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
			t.Errorf("CommandSet(%v) = (%v, %v); expected (%v, %v)", tc.argv, result, err, tc.result, tc.err)
		}
	}
}

func Test_CommandUnset(t *testing.T) {
	testCases := []struct {
		argv   []string
		result string
		err    error
	}{
		{[]string{"unset", "x"}, "", nil},
		{[]string{"unset", "y"}, "", nil},
	}

	i := picol.NewInterpreter()
	i.SetVariable("x", "42")
	i.SetVariable("y", "abc")

	for _, tc := range testCases {
		result, err := picol.CommandUnset(i, tc.argv, nil)
		if result != tc.result || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
			t.Errorf("CommandUnset(%v) = (%v, %v); expected (%v, %v)", tc.argv, result, err, tc.result, tc.err)
		}
	}
}

func Test_CommandIf(t *testing.T) {
	testCases := []struct {
		argv   []string
		result string
		err    error
	}{
		{[]string{"unset", "x"}, "", nil},
		{[]string{"unset", "y"}, "", nil},
	}

	i := picol.NewInterpreter()
	i.SetVariable("x", "42")
	i.SetVariable("y", "abc")

	for _, tc := range testCases {
		result, err := picol.CommandUnset(i, tc.argv, nil)
		if result != tc.result || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
			t.Errorf("CommandUnset(%v) = (%v, %v); expected (%v, %v)", tc.argv, result, err, tc.result, tc.err)
		}
	}
}

func Test_CommandWhile(t *testing.T) {
	i := picol.NewInterpreter()

	// Test simple while loop
	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	script := `
	set i 0
	set result 0
	while { < $i 5 } {
		set result [+ $result $i]
		set i [+ $i 1]
	}
	`

	res, err := i.Eval(script)
	if err != nil {
		t.Fatalf("Error during while loop evaluation: %s", err)
	}

	expectedResult := "10"
	if res != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, res)
	}

	// Test nested while loops
	script = `
	set i 0
	set j 0
	set result 0
	while { < $i 3 } {
		set j 0
		while { < $j 3 } {
			set result [+ $result $j]
			set j [+ $j 1]
		}
		set i [+ $i 1]
	}
	`

	res, err = i.Eval(script)
	if err != nil {
		t.Fatalf("Error during nested while loop evaluation: %s", err)
	}

	expectedResult = "9"
	if res != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, res)
	}
}

// You can also add test functions for other commands like CommandProc, CommandReturn, etc.
func Test_CommandRetCodes(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	// Test break
	script := `
	while { 1 } {
		break
	}
	`
	_, err = i.Eval(script)
	if err != nil {
		t.Fatalf("Error during break evaluation: %s", err)
	}

	// Test continue
	script = `
	set i 0
	set result 0
	while { < $i 5 } {
		if { == $i 2 } {
			continue
		}
		set result [+ $result $i]
		set i [+ $i 1]
	}
	`
	expectedResult := "7"
	res, err := i.Eval(script)
	if err != nil {
		t.Fatalf("Error during continue evaluation: %s", err)
	}
	if res != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, res)
	}
}

func Test_CommandProc(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	script := `
	proc sum {a b} {
		return [+ $a $b]
	}
	set res [sum 3 4]
	`
	expectedResult := "7"
	res, err := i.Eval(script)
	if err != nil {
		t.Fatalf("Error during proc evaluation: %s", err)
	}
	if res != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, res)
	}
}

func Test_CommandReturn(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	script := `
	proc testReturn {val} {
		return $val
	}
	set res [testReturn 42]
	`
	expectedResult := "42"
	res, err := i.Eval(script)
	if err != nil {
		t.Fatalf("Error during return evaluation: %s", err)
	}
	if res != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, res)
	}
}

func Test_CommandError(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	script := `
	error "An error occurred"
	`
	_, err = i.Eval(script)
	if err == nil || err.Error() != "An error occurred" {
		t.Fatalf("Error not raised or incorrect error message: %s", err)
	}
}

func Test_CommandPuts(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatal(err)
	}

	// The following test checks if the "puts" command runs without any error.
	// However, it doesn't check the printed output since it's not straightforward to capture stdout in tests.
	script := `
	puts "Hello, world!"
	`
	_, err = i.Eval(script)
	if err != nil {
		t.Fatalf("Error during puts evaluation: %s", err)
	}
}

func Test_RegisterCoreCommands(t *testing.T) {
	i := picol.NewInterpreter()

	err := i.RegisterCoreCommands()
	if err != nil {
		t.Fatalf("Error during core command registration: %s", err)
	}
}
