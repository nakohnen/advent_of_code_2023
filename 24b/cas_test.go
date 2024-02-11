package main

import (
	"testing"
)

func TestIntConstantEval(t *testing.T) {
	i := IntConstant(3)
	values := map[string]Expr{}
	result := i.Eval(values)
	if result != i {
		t.Errorf("Expected %v, got %v", i, result)
	}
}

func TestIntConstantEvalFloat(t *testing.T) {
	i := IntConstant(3)
	values := map[string]float64{}
	result := i.EvalFloat(values)
	if result != float64(i) {
		t.Errorf("Expected %v, got %v", 3, result)
	}
}

func TestIntConstantString(t *testing.T) {
	i := IntConstant(3)
	result := i.String()
	if result != "3" {
		t.Errorf("Expected %v, got %v", "3", result)
	}
}

func TestIntConstantExpand(t *testing.T) {
	i := IntConstant(3)
	result := i.Expand()
	if result != i {
		t.Errorf("Expected %v, got %v", i, result)
	}
}

func TestIntConstantSimplify(t *testing.T) {
	i := IntConstant(3)
	result := i.Simplify()
	if result != i {
		t.Errorf("Expected %v, got %v", i, result)
	}
}

func TestIntConstantResolveInner(t *testing.T) {
	i := IntConstant(3)
	result := i.ResolveInner()
	if result != i {
		t.Errorf("Expected %v, got %v", i, result)
	}
}

func TestFloatConstantEval(t *testing.T) {
	f := FloatConstant(3.0)
	values := map[string]Expr{}
	result := f.Eval(values)
	if result != f {
		t.Errorf("Expected %v, got %v", f, result)
	}
}

func TestFloatConstantEvalFloat(t *testing.T) {
	f := FloatConstant(3.0)
	values := map[string]float64{}
	result := f.EvalFloat(values)
	if result != float64(f) {
		t.Errorf("Expected %v, got %v", 3.0, result)
	}
}

func TestFloatConstantString(t *testing.T) {
	f := FloatConstant(3.0)
	result := f.String()
	if result != "3.000000" {
		t.Errorf("Expected %v, got %v", "3.000000", result)
	}
}

func TestFloatConstantExpand(t *testing.T) {
	f := FloatConstant(3.0)
	result := f.Expand()
	if result != f {
		t.Errorf("Expected %v, got %v", f, result)
	}
}

func TestFloatConstantSimplify(t *testing.T) {
	f := FloatConstant(3.0)
	result := f.Simplify()
	if result != f {
		t.Errorf("Expected %v, got %v", f, result)
	}
}

func TestFloatConstantResolveInner(t *testing.T) {
	f := FloatConstant(3.0)
	result := f.ResolveInner()
	if result != f {
		t.Errorf("Expected %v, got %v", f, result)
	}
}

func TestVariableEval(t *testing.T) {
	v := Variable("x")
	values := map[string]Expr{"x": IntConstant(3)}
	result := v.Eval(values)
	if result != IntConstant(3) {
		t.Errorf("Expected %v, got %v", IntConstant(3), result)
	}
}

func TestVariableEval2(t *testing.T) {
	v := Variable("x")
	values := map[string]Expr{"x": Variable("y")}
	result := v.Eval(values)
	if result != Variable("y") {
		t.Errorf("Expected %v, got %v", Variable("y"), result)
	}
}

func TestVariableEvalFloat(t *testing.T) {
	v := Variable("x")
	values := map[string]float64{"x": 3.0}
	result := v.EvalFloat(values)
	if result != 3.0 {
		t.Errorf("Expected %v, got %v", 3.0, result)
	}
}

func TestVariableString(t *testing.T) {
	v := Variable("x")
	result := v.String()
	if result != "x" {
		t.Errorf("Expected %v, got %v", "x", result)
	}
}

func TestVariableExpand(t *testing.T) {
	v := Variable("x")
	result := v.Expand()
	if result != v {
		t.Errorf("Expected %v, got %v", v, result)
	}
}

func TestVariableSimplify(t *testing.T) {
	v := Variable("x")
	result := v.Simplify()
	if result != v {
		t.Errorf("Expected %v, got %v", v, result)
	}
}

func TestVariableResolveInner(t *testing.T) {
	v := Variable("x")
	result := v.ResolveInner()
	if result != v {
		t.Errorf("Expected %v, got %v", v, result)
	}
}

func TestAdditionEval(t *testing.T) {
	a := Addition{[]Expr{IntConstant(3), IntConstant(4)}}
	values := map[string]Expr{}
	result := a.Eval(values)
	if result != IntConstant(7) {
		t.Errorf("Expected %v, got %v", IntConstant(7), result)
	}
}

func TestAdditionEval2(t *testing.T) {
	a := Addition{[]Expr{IntConstant(3), Variable("x")}}
	values := map[string]Expr{"x": IntConstant(4)}
	result := a.Eval(values)
	if result != IntConstant(7) {
		t.Errorf("Expected %v, got %v", IntConstant(7), result)
	}
}

func TestAdditionEvalFloat(t *testing.T) {
	a := Addition{[]Expr{IntConstant(3), IntConstant(4)}}
	values := map[string]float64{}
	result := a.EvalFloat(values)
	if result != 7.0 {
		t.Errorf("Expected %v, got %v", 7.0, result)
	}
}

func TestAdditionString(t *testing.T) {
	a := Addition{[]Expr{IntConstant(3), IntConstant(4)}}
	result := a.String()
	if result != "(3 + 4)" {
		t.Errorf("Expected %v, got %v", "(3 + 4)", result)
	}
}
