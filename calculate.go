/*
@Time : 2019/10/19 10:19
@Author : Hermes
@File : calculatePostfix
@Description:
*/
package calculator

import (
	"fmt"
	"strings"
)

func Calc(express string) (float64, error) {
	express = strings.ReplaceAll(express, " ", "")
	minors, err := calcExpressSplit(express, SUPPORT_SIGNS)
	if err != nil {
		return 0, err
	}
	minors = infix2ToPostfix(minors)
	return calculatePostfix(minors)
}

func calculatePostfix(postfix []*calcMinor) (float64, error) {
	stack := calcItemStack{}
	fixLen := len(postfix)
	for i := 0; i < fixLen; i++ {
		minor := postfix[i]
		// 数字：直接压栈
		if minor.Vtype == TYPE_NUMBER {
			stack.Push(minor)
		} else {
			// 操作符：取出两个数字计算值，再将结果压栈
			item := stack.Pop()
			if item == nil {
				return 0, fmt.Errorf("计算出错")
			}
			inval := item.Val
			fval2, ok := inval.(float64)
			if !ok {
				return 0, fmt.Errorf("不是合法的数值，%v", inval)
			}

			item = stack.Pop()
			if item == nil {
				return 0, fmt.Errorf("计算出错")
			}
			inval = item.Val
			fval1, ok := inval.(float64)
			if !ok {
				return 0, fmt.Errorf("不是合法的数值，%v", inval)
			}
			switch minor.Val {
			case "+":
				v := fval1 + fval2
				stack.Push(newCalcMinorItem(TYPE_NUMBER, v))
				//fmt.Printf("exec %v+%v=%v\r\n", fval1, fval2, v)
			case "-":
				v := fval1 - fval2
				stack.Push(newCalcMinorItem(TYPE_NUMBER, v))
				//fmt.Printf("exec %v-%v=%v\r\n", fval1, fval2, v)
			case "*":
				v := fval1 * fval2
				stack.Push(newCalcMinorItem(TYPE_NUMBER, v))
				//fmt.Printf("exec %v*%v=%v\r\n", fval1, fval2, v)
			case "/":
				v := fval1 / fval2
				stack.Push(newCalcMinorItem(TYPE_NUMBER, v))
				//fmt.Printf("exec %v/%v=%v\r\n", fval1, fval2, v)
			}
		}
	}

	item := stack.Pop()
	if item == nil {
		return 0, fmt.Errorf("计算出错")
	}
	inval := item.Val
	resVal, ok := inval.(float64)
	if !ok {
		return 0, fmt.Errorf("不是合法的数值，%v", inval)
	}
	return resVal, nil
}

func infix2ToPostfix(exp []*calcMinor) []*calcMinor {
	stack := calcItemStack{}
	var postfix []*calcMinor
	expLen := len(exp)

	// 遍历整个表达式
	for i := 0; i < expLen; i++ {

		minor := exp[i]
		if minor.Vtype == TYPE_NUMBER {
			postfix = append(postfix, minor)
		} else if minor.Vtype == TYPE_SIGN {
			switch minor.Val {
			case "(":
				// 左括号直接入栈
				stack.Push(minor)
			case ")":
				// 右括号则弹出元素直到遇到左括号
				for !stack.IsEmpty() {
					preminor := stack.Top()
					if preminor.Val == "(" {
						stack.Pop() // 弹出 "("
						break
					}
					postfix = append(postfix, preminor)
					stack.Pop()
				}
			default:
				// 操作符：遇到高优先级的运算符，不断弹出，直到遇见更低优先级运算符
				for !stack.IsEmpty() {
					top := stack.Top()
					if top.Val == "(" || isLower(top, minor) {
						break
					}
					postfix = append(postfix, top)
					stack.Pop()
				}
				// 低优先级的运算符入栈
				stack.Push(minor)
			}

		}
	}

	// 栈不空则全部输出
	for !stack.IsEmpty() {
		postfix = append(postfix, stack.Pop())
	}

	return postfix
}

func isLower(top *calcMinor, newTop *calcMinor) bool {
	// 注意 a + b + c 的后缀表达式是 ab + c +，不是 abc + +
	switch top.Val {
	case "+", "-":
		if newTop.Val == "*" || newTop.Val == "/" {
			return true
		}
	case "(":
		return true
	}
	return false
}
