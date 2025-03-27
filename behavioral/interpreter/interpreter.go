package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Context 上下文环境，用于存储变量和其对应的值
type Context struct {
	variables map[string]int
}

// NewContext 创建一个新的上下文环境
func NewContext() *Context {
	return &Context{
		variables: make(map[string]int),
	}
}

// SetVariable 设置变量值
func (c *Context) SetVariable(name string, value int) {
	c.variables[name] = value
}

// GetVariable 获取变量值
func (c *Context) GetVariable(name string) (int, bool) {
	value, exists := c.variables[name]
	return value, exists
}

// Expression 是解释器接口，定义了解释器的方法
type Expression interface {
	Interpret(context *Context) (int, error)
	String() string
}

// NumberExpression 表示一个数字字面量表达式
type NumberExpression struct {
	value int
}

// NewNumberExpression 创建一个数字表达式
func NewNumberExpression(value int) *NumberExpression {
	return &NumberExpression{value: value}
}

// Interpret 实现Expression接口，返回数字值
func (n *NumberExpression) Interpret(context *Context) (int, error) {
	return n.value, nil
}

// String 返回数字表达式的字符串表示
func (n *NumberExpression) String() string {
	return strconv.Itoa(n.value)
}

// VariableExpression 表示一个变量表达式
type VariableExpression struct {
	name string
}

// NewVariableExpression 创建一个变量表达式
func NewVariableExpression(name string) *VariableExpression {
	return &VariableExpression{name: name}
}

// Interpret 实现Expression接口，返回变量的值
func (v *VariableExpression) Interpret(context *Context) (int, error) {
	value, exists := context.GetVariable(v.name)
	if !exists {
		return 0, fmt.Errorf("变量 '%s' 未定义", v.name)
	}
	return value, nil
}

// String 返回变量表达式的字符串表示
func (v *VariableExpression) String() string {
	return v.name
}

// AddExpression 表示加法表达式
type AddExpression struct {
	left  Expression
	right Expression
}

// NewAddExpression 创建一个加法表达式
func NewAddExpression(left, right Expression) *AddExpression {
	return &AddExpression{left: left, right: right}
}

// Interpret 实现Expression接口，对左右表达式进行相加操作
func (a *AddExpression) Interpret(context *Context) (int, error) {
	leftValue, err := a.left.Interpret(context)
	if err != nil {
		return 0, err
	}

	rightValue, err := a.right.Interpret(context)
	if err != nil {
		return 0, err
	}

	return leftValue + rightValue, nil
}

// String 返回加法表达式的字符串表示
func (a *AddExpression) String() string {
	return fmt.Sprintf("(%s + %s)", a.left.String(), a.right.String())
}

// SubtractExpression 表示减法表达式
type SubtractExpression struct {
	left  Expression
	right Expression
}

// NewSubtractExpression 创建一个减法表达式
func NewSubtractExpression(left, right Expression) *SubtractExpression {
	return &SubtractExpression{left: left, right: right}
}

// Interpret 实现Expression接口，对左右表达式进行相减操作
func (s *SubtractExpression) Interpret(context *Context) (int, error) {
	leftValue, err := s.left.Interpret(context)
	if err != nil {
		return 0, err
	}

	rightValue, err := s.right.Interpret(context)
	if err != nil {
		return 0, err
	}

	return leftValue - rightValue, nil
}

// String 返回减法表达式的字符串表示
func (s *SubtractExpression) String() string {
	return fmt.Sprintf("(%s - %s)", s.left.String(), s.right.String())
}

// MultiplyExpression 表示乘法表达式
type MultiplyExpression struct {
	left  Expression
	right Expression
}

// NewMultiplyExpression 创建一个乘法表达式
func NewMultiplyExpression(left, right Expression) *MultiplyExpression {
	return &MultiplyExpression{left: left, right: right}
}

// Interpret 实现Expression接口，对左右表达式进行相乘操作
func (m *MultiplyExpression) Interpret(context *Context) (int, error) {
	leftValue, err := m.left.Interpret(context)
	if err != nil {
		return 0, err
	}

	rightValue, err := m.right.Interpret(context)
	if err != nil {
		return 0, err
	}

	return leftValue * rightValue, nil
}

// String 返回乘法表达式的字符串表示
func (m *MultiplyExpression) String() string {
	return fmt.Sprintf("(%s * %s)", m.left.String(), m.right.String())
}

// DivideExpression 表示除法表达式
type DivideExpression struct {
	left  Expression
	right Expression
}

// NewDivideExpression 创建一个除法表达式
func NewDivideExpression(left, right Expression) *DivideExpression {
	return &DivideExpression{left: left, right: right}
}

// Interpret 实现Expression接口，对左右表达式进行相除操作
func (d *DivideExpression) Interpret(context *Context) (int, error) {
	leftValue, err := d.left.Interpret(context)
	if err != nil {
		return 0, err
	}

	rightValue, err := d.right.Interpret(context)
	if err != nil {
		return 0, err
	}

	if rightValue == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}

	return leftValue / rightValue, nil
}

// String 返回除法表达式的字符串表示
func (d *DivideExpression) String() string {
	return fmt.Sprintf("(%s / %s)", d.left.String(), d.right.String())
}

// ModuloExpression 表示取模表达式
type ModuloExpression struct {
	left  Expression
	right Expression
}

// NewModuloExpression 创建一个取模表达式
func NewModuloExpression(left, right Expression) *ModuloExpression {
	return &ModuloExpression{left: left, right: right}
}

// Interpret 实现Expression接口，对左右表达式进行取模操作
func (m *ModuloExpression) Interpret(context *Context) (int, error) {
	leftValue, err := m.left.Interpret(context)
	if err != nil {
		return 0, err
	}

	rightValue, err := m.right.Interpret(context)
	if err != nil {
		return 0, err
	}

	if rightValue == 0 {
		return 0, fmt.Errorf("模数不能为零")
	}

	return leftValue % rightValue, nil
}

// String 返回取模表达式的字符串表示
func (m *ModuloExpression) String() string {
	return fmt.Sprintf("(%s %% %s)", m.left.String(), m.right.String())
}

// Parser 表达式解析器
type Parser struct {
	context *Context
	tokens  []string
	pos     int
}

// NewParser 创建一个新的解析器
func NewParser(context *Context) *Parser {
	return &Parser{
		context: context,
		tokens:  []string{},
		pos:     0,
	}
}

// Parse 解析表达式字符串并构建表达式树
func (p *Parser) Parse(expression string) (Expression, error) {
	// 词法分析，将表达式字符串拆分为标记
	p.tokenize(expression)
	p.pos = 0

	// 语法分析，构建表达式树
	return p.parseExpression()
}

// tokenize 将表达式字符串拆分为标记列表
func (p *Parser) tokenize(expression string) {
	p.tokens = []string{}

	// 去除所有空格
	expression = strings.ReplaceAll(expression, " ", "")

	i := 0
	for i < len(expression) {
		char := expression[i]

		// 处理数字
		if unicode.IsDigit(rune(char)) {
			num := ""
			for i < len(expression) && unicode.IsDigit(rune(expression[i])) {
				num += string(expression[i])
				i++
			}
			p.tokens = append(p.tokens, num)
			continue
		}

		// 处理变量名
		if unicode.IsLetter(rune(char)) {
			varName := ""
			for i < len(expression) && (unicode.IsLetter(rune(expression[i])) || unicode.IsDigit(rune(expression[i]))) {
				varName += string(expression[i])
				i++
			}
			p.tokens = append(p.tokens, varName)
			continue
		}

		// 处理运算符
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '%' || char == '(' || char == ')' {
			p.tokens = append(p.tokens, string(char))
			i++
			continue
		}

		// 跳过未知字符
		i++
	}
}

// parseExpression 解析加减表达式
func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]

		if token == "+" {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = NewAddExpression(left, right)
		} else if token == "-" {
			p.pos++
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = NewSubtractExpression(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parseTerm 解析乘除模表达式
func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]

		if token == "*" {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = NewMultiplyExpression(left, right)
		} else if token == "/" {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = NewDivideExpression(left, right)
		} else if token == "%" {
			p.pos++
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = NewModuloExpression(left, right)
		} else {
			break
		}
	}

	return left, nil
}

// parseFactor 解析因子（数字、变量、括号表达式）
func (p *Parser) parseFactor() (Expression, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("表达式意外结束")
	}

	token := p.tokens[p.pos]
	p.pos++

	// 处理括号表达式
	if token == "(" {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("缺少右括号")
		}
		p.pos++ // 跳过右括号
		return expr, nil
	}

	// 处理数字
	if num, err := strconv.Atoi(token); err == nil {
		return NewNumberExpression(num), nil
	}

	// 处理变量
	return NewVariableExpression(token), nil
}

// Evaluate 评估表达式字符串并返回结果
func Evaluate(expression string, context *Context) (int, error) {
	parser := NewParser(context)
	expr, err := parser.Parse(expression)
	if err != nil {
		return 0, err
	}

	return expr.Interpret(context)
}
