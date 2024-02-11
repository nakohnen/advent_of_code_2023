package main

import (
	"fmt"
	"sort"
	"strings"
)

type Expr interface {
	Eval(values map[string]Expr) Expr
	EvalFloat(values map[string]float64) float64
	Expand() Expr
	Simplify() Expr
	String() string
	ResolveInner() Expr
}

func removeElement(slice []Expr, elem Expr) []Expr {
	for i, v := range slice {
		if v == elem {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func RemoveDuplicatesStringSlice(slice [][]string) [][]string {
	seen := make(map[string]bool)
	result := [][]string{}
	for _, i := range slice {
		if !seen[strings.Join(i, "")] {
			result = append(result, i)
			seen[strings.Join(i, "")] = true
		}
	}
	return result
}

type IntConstant int

func (i IntConstant) Eval(values map[string]Expr) Expr {
	return i
}
func (i IntConstant) EvalFloat(values map[string]float64) float64 {
	return float64(i)
}
func (i IntConstant) String() string {
	return fmt.Sprintf("%d", i)
}
func (i IntConstant) Expand() Expr {
	return i
}
func (i IntConstant) Simplify() Expr {
	return i
}
func (i IntConstant) ResolveInner() Expr {
	return i
}

type FloatConstant float64

func (f FloatConstant) Eval(values map[string]Expr) Expr {
	return f
}
func (f FloatConstant) EvalFloat(values map[string]float64) float64 {
	return float64(f)
}
func (f FloatConstant) String() string {
	return fmt.Sprintf("%f", f)
}
func (f FloatConstant) Expand() Expr {
	return f
}
func (f FloatConstant) Simplify() Expr {
	return f
}
func (f FloatConstant) ResolveInner() Expr {
	return f
}

type Variable string

func (v Variable) Eval(values map[string]Expr) Expr {
	return values[string(v)]
}
func (v Variable) EvalFloat(values map[string]float64) float64 {
	return values[string(v)]
}
func (v Variable) String() string {
	return string(v)
}
func (v Variable) Expand() Expr {
	return v
}
func (v Variable) Simplify() Expr {
	return v
}
func (v Variable) ResolveInner() Expr {
	return v
}

type Addition struct {
	terms []Expr
}

func (a Addition) Eval(values map[string]Expr) Expr {
	b := a
	for i := 0; i < len(a.terms); i++ {
		b.terms[i] = a.terms[i].ResolveInner().Eval(values)
	}
	return b.ResolveInner()
}
func (a Addition) EvalFloat(values map[string]float64) float64 {
	sum := 0.0
	for _, term := range a.terms {
		sum += term.EvalFloat(values)
	}
	return sum
}
func (a Addition) String() string {
	s := "("
	for i, term := range a.terms {
		if i > 0 {
			s += " + "
		}
		s += term.String()
	}
	return s + ")"
}
func (a Addition) ResolveInner() Expr {
	//fmt.Println("Resolving inner", a.String())
	for i := 0; i < len(a.terms); i++ {
		//fmt.Println("->\t", a.terms[i].String())
		a.terms[i] = a.terms[i].ResolveInner()
	}

	// If sub-expression is Addition, flatten it
	//fmt.Println("If sub-expression is Addition, flatten it", a.String())
	new_terms := []Expr{}
	for _, term := range a.terms {
		if v, ok := term.(Addition); ok {
			inner_terms := []Expr{}
			for _, inner_term := range v.terms {
				inner_terms = append(inner_terms, inner_term.ResolveInner())
			}
			new_terms = append(new_terms, inner_terms...)
		} else {
			new_terms = append(new_terms, term.ResolveInner())
		}
	}

	// Sum all the constants of same type
	//fmt.Println("Sum all the constants of same type", new_terms)
	inner_int := IntConstant(0)
	other := []Expr{}
	for i := 0; i < len(new_terms); i++ {
		if v, ok := new_terms[i].(IntConstant); ok {
			inner_int += v
		} else {
			other = append(other, new_terms[i])
		}
	}

	inner_float := FloatConstant(0.0)
	other_tmp := []Expr{}
	for i := 0; i < len(other); i++ {
		if v, ok := other[i].(FloatConstant); ok {
			inner_float += v
		} else {
			other_tmp = append(other_tmp, other[i])
		}
	}
	//fmt.Println("Inner int", inner_int, "Inner float", inner_float, "Other", other_tmp)

	// Add all remaining terms which are symbols
	//fmt.Println("Add all remaining terms which are symbols", other_tmp)
	symbols := [][]string{}
	for _, term := range other_tmp {
		if v, ok := term.(Variable); ok {
			symbols = append(symbols, []string{v.String()})
		}
		if v, ok := term.(Multiplication); ok {
			symbols = append(symbols, v.GetSymbols())
		}
	}
	symbols = RemoveDuplicatesStringSlice(symbols)

	// Iterate over symbols and add their coefficients
	//fmt.Println("Iterate over symbols and add their coefficients", symbols)
	for _, symbol := range symbols {
		if len(symbol) == 0 {
			continue
		}
		symbol_repr := strings.Join(symbol, "*")
		terms := []Expr{}
		new_other_tmp := []Expr{}
		//fmt.Printf("Symbol \"%v\" for terms %v\n", symbol_repr, other_tmp)
		for _, term := range other_tmp {
			if v, ok := term.(Variable); ok {
				// Found "raw" symbol
				//fmt.Println("Found raw symbol")
				if v.String() == symbol_repr {
					terms = append(terms, IntConstant(1))
				} else {
					new_other_tmp = append(new_other_tmp, term)
				}
			} else if v, ok := term.(Multiplication); ok {
                // Found multiplication
                inner_symbols := v.GetSymbols()
                inner_symbols_repr := strings.Join(inner_symbols, "*")
				if inner_symbols_repr == symbol_repr {
					// Multiplication contains target symbols
					// Strip symbols from multiplication
					//fmt.Println("Multiplication contains target symbols")
					inner_terms := []Expr{}
					for _, inner_term := range v.terms {
						if _, ok := inner_term.(Variable); !ok {
							//fmt.Println("Inner term", inner_term.String(), "is not a variable")
							inner_terms = append(inner_terms, inner_term)
						}
					}
                    if len(inner_terms) == 0 {
                        terms = append(terms, IntConstant(1))
                    } else {
                        terms = append(terms, Multiplication{inner_terms})
                    }
				} else {
					new_other_tmp = append(new_other_tmp, term)
				}
			} else {
				new_other_tmp = append(new_other_tmp, term)
			}
		}
		new_coeff := Addition{terms}.ResolveInner()
		new_symbols := []Expr{}
		for _, s := range symbol {
			new_symbols = append(new_symbols, Variable(s))
		}
		new_term := Multiplication{append([]Expr{new_coeff}, new_symbols...)}.ResolveInner()

		// Remove old terms
		//fmt.Println("Remove old terms")
		other_tmp = new_other_tmp
		if new_coeff != IntConstant(0) && new_coeff != FloatConstant(0.0) {
			other_tmp = append(other_tmp, new_term)
		}
	}

	if inner_int != 0 {
		other_tmp = append(other_tmp, inner_int)
	}
	if inner_float != 0.0 {
		other_tmp = append(other_tmp, inner_float)
	}
	if len(other_tmp) == 1 {
		return other_tmp[0]
	} else if len(other_tmp) == 0 {
		return IntConstant(0)
	}
	return Addition{other_tmp}
}
func (a Addition) Expand() Expr {
	if len(a.terms) == 1 {
		return a.terms[0].Expand()
	}
	for i := 0; i < len(a.terms); i++ {
		a.terms[i] = a.terms[i].Expand()
	}
	return a
}
func (a Addition) Simplify() Expr {
	return a
}

type Multiplication struct {
	terms []Expr
}

func (m Multiplication) Eval(values map[string]Expr) Expr {
	b := m
	for i := 0; i < len(m.terms); i++ {
		b.terms[i] = m.terms[i].Eval(values)
	}
	return b
}
func (m Multiplication) EvalFloat(values map[string]float64) float64 {
	product := 1.0
	for _, term := range m.terms {
		product *= term.EvalFloat(values)
	}
	return product
}
func (m Multiplication) String() string {
	s := ""
	for i, term := range m.terms {
		if i > 0 {
			s += "*"
		}
		s += term.String()
	}
	return s
}
func (m Multiplication) Expand() Expr {
	m.ResolveInner()
	expanded := false
	if len(m.terms) == 1 {
		return m.terms[0].Expand()
	}
	for i := 0; i < len(m.terms); i++ {
		new_term := m.terms[i].Expand()
		if new_term.String() != m.terms[i].String() {
			expanded = true
		}
		m.terms[i] = m.terms[i].Expand()
	}
	m.ResolveInner()
	found_addition := 0
	found_int_constant := 0
	found_float_constant := 0
	found_variable := 0
	found_fraction := 0
	for _, term := range m.terms {
		if _, ok := term.(Addition); ok {
			found_addition++
		}
		if _, ok := term.(IntConstant); ok {
			found_int_constant++
		}
		if _, ok := term.(FloatConstant); ok {
			found_float_constant++
		}
		if _, ok := term.(Variable); ok {
			found_variable++
		}
		if _, ok := term.(Fraction); ok {
			found_fraction++
		}
	}
	if found_int_constant > 0 && (found_addition > 0 || found_fraction > 0) {
		int_constant := IntConstant(0)
		for i, term := range m.terms {
			if v, ok := term.(IntConstant); ok {
				int_constant = v
				m.terms = append(m.terms[:i], m.terms[i+1:]...)
				break
			}
		}
		if found_addition > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Addition); ok {
					for j, inner_term := range v.terms {
						new_term := Multiplication{[]Expr{int_constant, inner_term}}.ResolveInner()
						v.terms[j] = new_term
						expanded = true
					}
					break
				}
			}
		} else if found_fraction > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Fraction); ok {
					v.num = Multiplication{[]Expr{int_constant, v.num}}.ResolveInner()
					expanded = true
					break
				}
			}
		}
	}
	if found_float_constant > 0 && (found_addition > 0 || found_fraction > 0) {
		float_constant := FloatConstant(0.0)
		for i, term := range m.terms {
			if v, ok := term.(FloatConstant); ok {
				float_constant = v
				m.terms = append(m.terms[:i], m.terms[i+1:]...)
				expanded = true
				break
			}
		}
		if found_addition > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Addition); ok {
					for j, inner_term := range v.terms {
						v.terms[j] = Multiplication{[]Expr{float_constant, inner_term}}.ResolveInner()
					}
					expanded = true
					break
				}
			}
		} else if found_fraction > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Fraction); ok {
					v.num = Multiplication{[]Expr{float_constant, v.num}}.ResolveInner()
					expanded = true
					break
				}
			}
		}

	}
	if found_variable == 1 && (found_addition > 0 || found_fraction > 0) {
		sym := Variable("")
		for i, term := range m.terms {
			if v, ok := term.(Variable); ok {
				sym = v
				m.terms = append(m.terms[:i], m.terms[i+1:]...)
				//fmt.Println("Found symbol", sym)
				break
			}
		}
		if sym.String() == "" {
			panic(fmt.Sprintf("Variable not found in %v", m))
		}
		if found_addition > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Addition); ok {
					for j, inner_term := range v.terms {
						v.terms[j] = Multiplication{[]Expr{sym, inner_term}}.ResolveInner()
					}
					expanded = true
					break
				}
			}
		} else if found_fraction > 0 {
			// Distribute
			for _, term := range m.terms {
				if v, ok := term.(Fraction); ok {
					v.num = Multiplication{[]Expr{sym, v.num}}.ResolveInner()
					expanded = true
					break
				}
			}
		}
	} else if found_addition > 1 {
		// Distribute
		add1 := Addition{[]Expr{}}
		add1_index := -1
		for i, term := range m.terms {
			if v, ok := term.(Addition); ok {
				add1 = v
				add1_index = i
				break
			}
		}
		m.terms = append(m.terms[:add1_index], m.terms[add1_index+1:]...)
		for i := 0; i < len(m.terms); i++ {
			term := m.terms[i]
			if v, ok := term.(Addition); ok {
				// Distribute inner
				new_terms := []Expr{}
				for _, inner_term := range v.terms {
					for _, inner_add_term := range add1.terms {
						new_terms = append(new_terms, Multiplication{[]Expr{inner_term, inner_add_term}}.ResolveInner())
					}
				}
				m.terms[i] = Addition{new_terms}.ResolveInner()
				expanded = true
				break
			}
		}
	}
	if expanded {
		//fmt.Println("Expanded", m.String())
		return m.Expand()
	}
	return m.ResolveInner()
}
func (m Multiplication) Simplify() Expr {
	return m
}
func (m Multiplication) GetSymbols() []string {
	symbols := []string{}
	for _, term := range m.terms {
		if v, ok := term.(Variable); ok {
			symbols = append(symbols, string(v))
		}
	}
	return symbols
}
func (m Multiplication) ResolveInner() Expr {
	for i := 0; i < len(m.terms); i++ {
		m.terms[i] = m.terms[i].ResolveInner()
	}

	remove_one := false
	for remove_one {
		remove_one = false
		for i := 0; i < len(m.terms); i++ {
			if v, ok := m.terms[i].(IntConstant); ok {
				if v == 1 {
					m.terms = append(m.terms[:i], m.terms[i+1:]...)
					remove_one = true
					break
				}
			} else if v, ok := m.terms[i].(FloatConstant); ok {
				if v == 1.0 {
					m.terms = append(m.terms[:i], m.terms[i+1:]...)
					remove_one = true
					break
				}
			}
		}
	}

	// If sub-expression is Multiplication, flatten it
	new_terms := []Expr{}
	for _, term := range m.terms {
		if v, ok := term.(Multiplication); ok {
			inner_terms := []Expr{}
			for _, inner_term := range v.terms {
				inner_terms = append(inner_terms, inner_term.ResolveInner())
			}
			new_terms = append(new_terms, inner_terms...)
		} else {
			new_terms = append(new_terms, term)
		}
	}

    // If sub-expression is Fraction, flatten it
    tmp_terms := []Expr{}
    fractions := []Expr{}
    other_terms := []Expr{}
    for _, term := range new_terms {
        if _, ok := term.(Fraction); ok {
            fractions = append(fractions, term)
        } else {
            other_terms = append(other_terms, term)
        }
    }
    new_num := []Expr{}
    new_den := []Expr{}
    for _, term := range fractions {
        new_num = append(new_num, term.(Fraction).num)
        new_den = append(new_den, term.(Fraction).den)
    }
    if len(new_num) > 0 {
        new_fraction := Fraction{Multiplication{new_num}, Multiplication{new_den}}
        tmp_terms = append(tmp_terms, new_fraction.ResolveInner())
    }
    tmp_terms = append(tmp_terms, other_terms...)
    new_terms = tmp_terms



	// Sort terms
	//fmt.Println("Before sort", new_terms)
	sort.Slice(new_terms, func(i, j int) bool {
		return new_terms[i].String() < new_terms[j].String()
	})
	//fmt.Println("After sort", new_terms)

	// Sum all the constants of same type
	inner_int := IntConstant(1)
	other := []Expr{}
	for i := 0; i < len(new_terms); i++ {
		if v, ok := new_terms[i].(IntConstant); ok {
			inner_int *= v
		} else {
			other = append(other, new_terms[i])
		}
	}
	if inner_int != 1 {
		// Add to front of list
		other = append([]Expr{IntConstant(inner_int)}, other...)
	}

	inner_float := FloatConstant(1.0)
	other_tmp := []Expr{}
	for i := 0; i < len(other); i++ {
		if v, ok := other[i].(FloatConstant); ok {
			inner_float *= v
		} else {
			other_tmp = append(other_tmp, other[i])
		}
	}
	if inner_float != 1.0 {
		// Add to front of list
		other_tmp = append([]Expr{FloatConstant(inner_float)}, other_tmp...)
	}
	if len(other_tmp) == 1 {
		return other_tmp[0]
	}

	// If there is a zero, return zero
	found_zero := false
	for _, term := range other_tmp {
		if v, ok := term.(IntConstant); ok {
			if v == 0 {
				found_zero = true
				break
			}
		}
		if v, ok := term.(FloatConstant); ok {
			if v == 0.0 {
				found_zero = true
				break
			}
		}
	}
	if found_zero {
		return IntConstant(0)
	}

	return Multiplication{other_tmp}
}

type Fraction struct {
	num, den Expr
}

func (f Fraction) Eval(values map[string]Expr) Expr {
	if f.den.Eval(values) == IntConstant(1) || f.den.Eval(values) == FloatConstant(1.0) {
		return f.num.Eval(values)
	}
	return Fraction{f.num.Eval(values), f.den.Eval(values)}
}
func (f Fraction) EvalFloat(values map[string]float64) float64 {
	return f.num.EvalFloat(values) / f.den.EvalFloat(values)
}
func (f Fraction) String() string {
	return f.num.String() + "/" + f.den.String()
}
func (f Fraction) Expand() Expr {
    f.num = f.num.Expand()
    f.den = f.den.Expand()
	return f.ResolveInner()
}
func (f Fraction) Simplify() Expr {
    f.num = f.num.Simplify()
    f.den = f.den.Simplify()
	return f.ResolveInner()
}
func (f Fraction) ResolveInner() Expr {
    // If numerator is 0, return 0
    num := f.num.ResolveInner()
    if v, ok := num.(IntConstant); ok {
        if v == 0 {
            return IntConstant(0)
        }
    }
    // If denominator is 1, return numerator
	den := f.den.ResolveInner()
	if v, ok := den.(IntConstant); ok {
		if v == 1 {
			return f.num.ResolveInner()
		}
	}
	if v, ok := den.(FloatConstant); ok {
		if v == 1.0 {
			return f.num.ResolveInner()
		}
	}
    
    // If numerator and denominator are the same, return 1
    if f.num.String() == f.den.String() {
        return IntConstant(1)
    }

	return Fraction{f.num.ResolveInner(), f.den.ResolveInner()}
}

type Equation struct {
	left, right Expr
}

func (e Equation) Eval(values map[string]Expr) Equation {
    return Equation{e.left.Eval(values), e.right.Eval(values)}
}

func (e Equation) String() string {
    return e.left.String() + " = " + e.right.String()
}

/* func main() {
	expressions := []Expr{}

	// (x +1) * (y + 2) * (z + 3)
	inner_sub20 := Addition{[]Expr{Variable("x"), IntConstant(1)}}
	inner_sub21 := Addition{[]Expr{Variable("y"), IntConstant(2)}}
	inner_sub22 := Addition{[]Expr{Variable("z"), IntConstant(3)}}
	expressions = append(expressions, Multiplication{[]Expr{inner_sub20, inner_sub21, inner_sub22}})

    // (x+1)/(x+1)
    inner_sub23 := Addition{[]Expr{Variable("x"), IntConstant(1)}}
    expressions = append(expressions, Fraction{inner_sub23, inner_sub23})

    // (x+1)/1
    expressions = append(expressions, Fraction{inner_sub23, IntConstant(1)})

    // (x+1)/(x+2) * (x+2)/(x+1)
    inner_sub24 := Addition{[]Expr{Variable("x"), IntConstant(2)}}
    inner_sub25 := Fraction{inner_sub23, inner_sub24}
    inner_sub26 := Fraction{inner_sub24, inner_sub23}
    expressions = append(expressions, Multiplication{[]Expr{inner_sub25, inner_sub26}})

    // (x + (-1)*x)/(x + 1)
    inner_sub27 := Multiplication{[]Expr{Addition{[]Expr{Variable("x"), Multiplication{[]Expr{IntConstant(-1), Variable("x")}}}}}}
    inner_sub28 := Addition{[]Expr{Variable("x"), IntConstant(1)}}
    expressions = append(expressions, Fraction{inner_sub27, inner_sub28})

    // (-2*x + 2*x)/(x + 1)
    inner_sub29 := Multiplication{[]Expr{IntConstant(-2), Variable("x")}}
    inner_sub30 := Multiplication{[]Expr{IntConstant(2), Variable("x")}}
    inner_sub31 := Addition{[]Expr{inner_sub29, inner_sub30}}
    inner_sub32 := Addition{[]Expr{Variable("x"), IntConstant(1)}}
    expressions = append(expressions, Fraction{inner_sub31, inner_sub32})

    // (x + 1)/(d + 2) * (y + z + 5)
    inner_sub33 := Addition{[]Expr{Variable("y"), Variable("z"), IntConstant(5)}}
    inner_sub34 := Addition{[]Expr{Variable("d"), IntConstant(2)}}
    inner_sub35 := Fraction{inner_sub23, inner_sub34}
    expressions = append(expressions, Multiplication{[]Expr{inner_sub35, inner_sub33}})

	for _, expr := range expressions {
		fmt.Println("Original", expr.String())
		fmt.Println("Expanded", expr.Expand().ResolveInner().String())
		fmt.Println()
	}

} */
