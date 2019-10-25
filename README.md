# go-simple-calculator
a simple calculator by golang support +-*/()

一个简单的纯go实现的计算器
支持普通四则运算和括号优先符号


能够处理的表达式样例：

    1+127-21+(3-4)*6/2.5
    
    (88+(1+8)*6)/2+99

直接调用解析执行函数 :

    func main() {

        s := "-11.1+2.22*(73.45/6.463)"
        val, err := calculator.Calc(express)
  
        if err != nil {
        
            return
            
        }
        
        fmt.Printf("计算结果: %f\n", val)
}