package rpn

import (
	"errors"
	"fmt"
	gtoken "go/token"
	"math"
	"strconv"
)

var (
	ErrConvertFuncIsNil       = errors.New("function for convert value to float can be nil only if values is nil")
	ErrComputationalOperation = errors.New("not computational parametres for computational operation")
	ErrComparisonOperation    = errors.New("not computational parametres for comparison operation")
	ErrLogicalOperation       = errors.New("not logical parametres for logical operation")
)

type executeTokenType int

const (
	FLOAT executeTokenType = iota
	BOOL
)

// executeToken - вспомогательный токен, чтобы не парсить постоянно значение из строки в строку.
// Тип int будет приведен к float64
type executeToken struct {
	kind  executeTokenType
	value float64
}

func (et *executeToken) IsBoolean() bool {
	if et.kind == BOOL {
		return true
	}
	return false
}

func (et *executeToken) IsFloat() bool {
	if et.kind == FLOAT {
		return true
	}
	return false
}

type toFloat64[T any] func(T) float64

// Execute исполняет rpn для набора токенов, values - содержит набор текстовых параметров, to - функция для преобразования найденого значение в values во float64.
// Если переданая rpn не содержит текстовых параметров, то values и to могут быть nil.
func Execute[T any](tokens []*Token, values map[string]T, to toFloat64[T]) (*executeToken, error) {
	stack := NewStack[*executeToken]()
	for _, token := range tokens {
		if token.Kind.IsLiteral() {
			switch token.Kind {
			case gtoken.INT, gtoken.FLOAT:
				f64v, err := strconv.ParseFloat(token.Value, 64)
				if err != nil {
					return nil, err
				}
				stack.Push(&executeToken{FLOAT, f64v})
			case gtoken.IDENT:
				value, ok := values[token.Value]
				if !ok {
					// TODO: Custom type err
					return nil, fmt.Errorf("not found value for '%v'", token.Value)
				}
				if to == nil {
					return nil, ErrConvertFuncIsNil
				}
				stack.Push(&executeToken{FLOAT, to(value)})
			}
		}
		if token.Kind.IsOperator() {
			v2 := stack.Pop()
			v1 := stack.Pop()
			if operation, ok := opertations[token.Kind]; ok {
				var tokenType executeTokenType
				switch getOperationType(token.Kind) {
				case COMPUTATIONAL:
					if !v1.IsFloat() || !v2.IsFloat() {
						return nil, ErrComputationalOperation
					}
					tokenType = FLOAT
				case COMPARISON:
					if !v1.IsFloat() || !v2.IsFloat() {
						return nil, ErrComparisonOperation
					}
					tokenType = BOOL
				case LOGIC:
					if !v1.IsBoolean() || !v2.IsBoolean() {
						return nil, ErrLogicalOperation
					}
					tokenType = BOOL
				}
				stack.Push(&executeToken{tokenType, operation(v1.value, v2.value)})
			}
		}
	}
	return stack.Pop(), nil
}

func boolToFloat64(condition bool) float64 {
	if condition {
		return 1
	}
	return 0
}

func float64ToBool(value float64) bool {
	if value == 0 {
		return false
	}
	return true
}

type operationType int

const (
	NOTSUPPORT operationType = iota
	COMPUTATIONAL
	COMPARISON
	LOGIC
)

func getOperationType(kind gtoken.Token) operationType {
	switch kind {
	case gtoken.ADD, gtoken.SUB, gtoken.MUL, gtoken.QUO, gtoken.XOR:
		return COMPUTATIONAL
	case gtoken.EQL, gtoken.LSS, gtoken.GTR, gtoken.NEQ, gtoken.LEQ, gtoken.GEQ:
		return COMPARISON
	case gtoken.OR, gtoken.AND:
		return LOGIC
	}
	return NOTSUPPORT
}

var opertations = map[gtoken.Token]func(a, b float64) float64{
	// computational
	gtoken.ADD: func(a, b float64) float64 { return a + b },          // +
	gtoken.SUB: func(a, b float64) float64 { return a - b },          // -
	gtoken.MUL: func(a, b float64) float64 { return a * b },          // *
	gtoken.QUO: func(a, b float64) float64 { return a / b },          // /
	gtoken.XOR: func(a, b float64) float64 { return math.Pow(a, b) }, // ^
	// comparison
	gtoken.EQL: func(a, b float64) float64 { return boolToFloat64(a == b) }, // ==
	gtoken.LSS: func(a, b float64) float64 { return boolToFloat64(a < b) },  // <
	gtoken.GTR: func(a, b float64) float64 { return boolToFloat64(a > b) },  // >
	gtoken.NEQ: func(a, b float64) float64 { return boolToFloat64(a != b) }, // !=
	gtoken.LEQ: func(a, b float64) float64 { return boolToFloat64(a <= b) }, // <=
	gtoken.GEQ: func(a, b float64) float64 { return boolToFloat64(a >= b) }, // >=
	// logic
	gtoken.OR:  func(a, b float64) float64 { return boolToFloat64(float64ToBool(a) || float64ToBool(b)) }, // OR
	gtoken.AND: func(a, b float64) float64 { return boolToFloat64(float64ToBool(a) && float64ToBool(b)) }, // AND
}
