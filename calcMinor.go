/*
@Time : 2019/10/19 12:41
@Author : Hermes
@File : calcMinor
@Description:
*/
package calculator

import (
	"bytes"
	"fmt"
)

const TYPE_NUMBER = 1
const TYPE_SIGN = 2

type calcMinor struct {
	Vtype int
	Val   interface{} //固定数值类型的时候是float64,符号类型的时候是string
}

func newCalcMinorItem(t int, val interface{}) *calcMinor {
	return &calcMinor{
		Vtype: t,
		Val:   val,
	}
}

func joinCalcMinors(arr []*calcMinor) string {
	bf := new(bytes.Buffer)
	for idx := range arr {
		bf.WriteString(fmt.Sprintf("%v", arr[idx].Val))
	}
	return bf.String()
}
