package generate

import "github.com/ivan-jorge001/gormt/tools"

// interval.间隔
var _interval = "\t"

// IGenerate Generate Printing Interface.生成打印接口
type IGenerate interface {
	Generate() string
}

// PrintAtom . atom print .原始打印
type PrintAtom struct {
	lines []string
}

// Add  one to print.打印
func (p *PrintAtom) Add(str ...interface{}) {
	var tmp string
	for _, v := range str {
		tmp += tools.AsString(v) + _interval
	}
	p.lines = append(p.lines, tmp)
}

// Generates Get the generated list.获取生成列表
func (p *PrintAtom) Generates() []string {
	return p.lines
}
