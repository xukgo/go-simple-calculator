/*
@Time : 2019/10/19 10:25
@Author : Hermes
@File : splitExpress
@Description:
*/
package calculator

import (
	"bytes"
	"fmt"
	"github.com/xukgo/gsaber/utils/ruleUtil"
	"github.com/xukgo/gsaber/utils/stringUtil"
	"strconv"
)

//把表达式分解成数字和符号，过程中自动校验数字是否是浮点数，自动调整一些符号比如开头和括号后开头的+-
func calcExpressSplit(exp string, signs []string) ([]*calcMinor, error) {
	var err error
	var arr []*calcMinor
	srclen := len(exp)

	leftBracketCount := 0
	rightBracketCount := 0
	var first bool = true
	index := 0     //代表src指针目前的偏移
	lastSign := "" //代表上个是符号，不是则置为空
	bf := new(bytes.Buffer)
	for i := 0; i < srclen; i++ {
		signLen := 0
		matchSign := ""
		for idx := range signs {
			if stringUtil.StartWithIndex(exp, index, signs[idx]) {
				matchSign = signs[idx]
				signLen = len(matchSign)
				break
			}
		}

		if signLen == 0 {
			bf.WriteByte(exp[i])
			index++
			lastSign = ""
		} else {
			if bf.Len() > 0 {
				numberVal := bf.String()
				if !ruleUtil.CheckIsFloat(numberVal) {
					return nil, fmt.Errorf("表达式内不是有效的数值类型:%s", numberVal)
				}
				fval, err := strconv.ParseFloat(numberVal, 64)
				if err != nil {
					return nil, fmt.Errorf("表达式内不是有效的数值类型:%s", numberVal)
				}
				arr = append(arr, newCalcMinorItem(TYPE_NUMBER, fval))
				bf.Reset()
			}

			arr, err = handlerSign(arr, bf, first, lastSign, matchSign)
			if err != nil {
				return nil, err
			}
			if matchSign == "(" {
				leftBracketCount++
			} else if matchSign == ")" {
				rightBracketCount++
			}
			index += signLen
			lastSign = matchSign
		}

		first = false
	}

	//结尾后处理，如果是数字还要加进去
	if bf.Len() > 0 {
		numberVal := bf.String()
		if !ruleUtil.CheckIsFloat(numberVal) {
			return nil, fmt.Errorf("表达式内不是有效的数值类型:%s", numberVal)
		}
		fval, err := strconv.ParseFloat(numberVal, 64)
		if err != nil {
			return nil, fmt.Errorf("表达式内不是有效的数值类型:%s", numberVal)
		}
		arr = append(arr, newCalcMinorItem(TYPE_NUMBER, fval))
		//bf.Reset()
	} else {
		if len(arr) > 0 {
			//验证最后一个如果是符号，是否合法,最后结尾只允许)
			lastMinor := arr[len(arr)-1]
			if lastMinor.Val != ")" {
				return nil, fmt.Errorf("表达式最后结尾作为符号只允许)")
			}
		}
	}
	//验证()数量对应
	if leftBracketCount != rightBracketCount {
		return nil, fmt.Errorf("左括号和右括号的数量不匹配")
	}
	return arr, nil
}

func handlerSign(arr []*calcMinor, bf *bytes.Buffer, first bool, lastSign string, matchSign string) ([]*calcMinor, error) {
	if matchSign == "+" {
		//这些情况忽略+
		if first || lastSign == "*" || lastSign == "/" || lastSign == "(" {
			//忽略
		} else if lastSign == "+" || lastSign == "-" { //这些情况不允许，报错
			return nil, fmt.Errorf("请检查表达式，不允许+-叠加在一起")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	} else if matchSign == "-" {
		//这些情况入数值缓存
		if first || lastSign == "*" || lastSign == "/" || lastSign == "(" {
			bf.WriteString(matchSign)
		} else if lastSign == "+" || lastSign == "-" { //这些情况不允许，报错
			return nil, fmt.Errorf("请检查表达式，不允许+-叠加在一起")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	} else if matchSign == "*" {
		if first {
			return nil, fmt.Errorf("请检查表达式，不允许以*开始")
		} else if lastSign == "*" || lastSign == "/" || lastSign == "+" || lastSign == "-" || lastSign == "(" {
			return nil, fmt.Errorf("请检查表达式，不允许非法符号叠加")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	} else if matchSign == "/" {
		if first {
			return nil, fmt.Errorf("请检查表达式，不允许以/开始")
		} else if lastSign == "*" || lastSign == "/" || lastSign == "+" || lastSign == "-" || lastSign == "(" {
			return nil, fmt.Errorf("请检查表达式，不允许非法符号叠加")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	} else if matchSign == "(" {
		if lastSign == ")" {
			return nil, fmt.Errorf("请检查表达式，不允许非法符号叠加")
		} else if !first && lastSign == "" {
			return nil, fmt.Errorf("请检查表达式，不允许数值后面直接接(")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	} else if matchSign == ")" {
		if first {
			return nil, fmt.Errorf("请检查表达式，不允许符号)开始")
		} else if lastSign == "(" || lastSign == "+" || lastSign == "-" || lastSign == "*" || lastSign == "/" {
			return nil, fmt.Errorf("请检查表达式，不允许非法符号叠加")
		} else { // ")"正常入队
			arr = append(arr, newCalcMinorItem(TYPE_SIGN, matchSign))
		}
	}

	return arr, nil
}

//
////验证和调整表达式分解后的顺序合法
//func checkAndAdaptExpressSplits(minors []*calcMinor) ([]*calcMinor, error) {
//	var delIndexArr []int
//	var lastSign string
//	var currentSign string
//	if checkIsSign(minors[0], "+") {
//		delIndexArr = append(delIndexArr, 0)
//	}
//	if checkIsSign(minors[0], "-") {
//		delIndexArr = append(delIndexArr, 0)
//	}
//	return minors, nil
//}

func checkIsSign(minor *calcMinor, sign string) bool {
	if minor.Vtype != TYPE_SIGN {
		return false
	}
	if minor.Val == sign {
		return true
	}

	return false
}

func removeCalcMinorIndex(arr []*calcMinor, index int) []*calcMinor {
	j := 0
	for idx := range arr {
		if idx != index {
			arr[j] = arr[idx]
			j++
		}
	}
	return arr[:j]
}
func insertCalcMinorIndex(arr []*calcMinor, index int, addItem *calcMinor) []*calcMinor {
	j := 0
	for idx := range arr {
		if idx != index {
			arr[j] = arr[idx]
			j++
		}
	}
	return arr[:j]
}
