package evaluator

import (
	"eventloop/interpreter/interpreter/lexer"
	"eventloop/interpreter/interpreter/object"
	"eventloop/interpreter/interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"(1 <= 2) == true", true},
		{"(2 >= 1) == true", true},
		{"(1 >= 2) == true", false},
		{"(2 <= 1) == true", false},
		{"true && true", true},
		{"1 && 2", true},
		{"false || true", true},
		{"0 || 1", true},
		{"[] || 1", true},
		{`[] || "true"`, true},
		{`"" || "true"`, true},
		{`"" || [true]`, true},
		{`"" || [1, 2]`, true},
		{"false && false", false},
		{"true && false", false},
		{"false && true", false},
		{"0 && 0", false},
		{"0 || 0", false},
		{"false || false", false},
		{`"" && ""`, false},
		{`"" && []`, false},
		{`"" && true`, false},
		{`[] && true`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
		  if (10 > 1) {
			return 10;
		  }

		  return 1;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
			  if (10 > 1) {
				return true + false;
			  }

			  return 1;
			}
		`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"x = 31 + 38",
			"identifier not found: x",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestReAssignmentStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 5; a = a + 1; a;", 6},
		{`let a = "lol"; a = [1, 2, 3]; a[0];`, 1},
		{`let s = "Hello "; s = s + "World!"; s;`, "Hello World!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(tt.expected.(int)))
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}

			if str.Value != tt.expected.(string) {
				t.Errorf("String has wrong value. got=%q", str.Value)
			}
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
	let first = 10;
	let second = 10;
	let third = 10;

	let ourFunction = fn(first) {
	  let second = 20;

	  first + second + third;
	};

	ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
	  fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`rest([1, 2, 3])`, &object.Array{Elements: []object.Object{&object.Integer{Value: 2}, &object.Integer{Value: 3}}}},
		{`rest([1])`, &object.Array{Elements: []object.Object{}}},
		{`rest([])`, nil},
		{`push([1], 2)`, &object.Array{Elements: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}}}},
		{`push([], 1)`, &object.Array{Elements: []object.Object{&object.Integer{Value: 1}}}},
		{`pop([1, 2, 3])`, &object.Array{Elements: []object.Object{&object.Integer{Value: 1}, &object.Integer{Value: 2}}}},
		{`pop([1])`, &object.Array{Elements: []object.Object{}}},
		{`pop([])`, nil},
		{`bool(true)`, true},
		{`bool(false)`, false},
		{`bool("")`, false},
		{`bool([])`, false},
		{`bool("false")`, true},
		{`bool([false])`, true},
		{`print("")`, nil},
		{`print()`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(tt.expected.(int)))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case nil:
			testNullObject(t, evaluated)
		case *object.Array:
			expected_arr, _ := tt.expected.(*object.Array)
			testArrayObject(t, evaluated, *expected_arr)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testArrayObject(t *testing.T, obj object.Object, expected object.Array) bool {
	arr, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not ARRAY. got=%T (%+v)", obj, obj)
	}

	len_expected := len(expected.Elements)
	len_arr := len(arr.Elements)
	if len_expected != len_arr {
		t.Errorf("expected array of size=%d got=%d", len_expected, len_arr)
	}

	for i, expected_elem := range expected.Elements {
		curr_elem := arr.Elements[i]
		switch expected_elem.(type) {
		case *object.Integer:
			expected_int, _ := expected_elem.(*object.Integer)
			testIntegerObject(t, curr_elem, expected_int.Value)
		case *object.Boolean:
			expected_bool, _ := expected_elem.(*object.Boolean)
			testBooleanObject(t, curr_elem, expected_bool.Value)
		case *object.Null:
			testNullObject(t, curr_elem)
		}
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}
