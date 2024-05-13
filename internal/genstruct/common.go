package genstruct

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wonli/gormt/internal/cnf"
	"github.com/wonli/gormt/internal/generate"
)

// SetName Setting element name.设置元素名字
func (e *GenElement) SetName(name string) {
	e.Name = name
}

// SetType Setting element type.设置元素类型
func (e *GenElement) SetType(tp string) {
	e.Type = tp
}

// SetNotes Setting element notes.设置注释
func (e *GenElement) SetNotes(notes string) {
	e.Notes = strings.Replace(notes, "\n", ",", -1)
}

// AddTag Add a tag .添加一个tag标记
func (e *GenElement) AddTag(k string, v string) {
	if e.Tags == nil {
		e.Tags = make(map[string][]string)
	}
	e.Tags[k] = append(e.Tags[k], v)
}

// Generate Get the result data.获取结果数据
func (e *GenElement) Generate() string {
	tag := ""
	if e.Tags != nil {
		var ks []string
		for k := range e.Tags {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		var tags []string
		for _, v := range ks {
			tags = append(tags, fmt.Sprintf(`%v:"%v"`, v, strings.Join(e.Tags[v], ";")))
		}
		tag = fmt.Sprintf("`%v`", strings.Join(tags, " "))
	}

	var p generate.PrintAtom
	if len(e.Notes) > 0 {
		p.Add(e.Name, e.Type, tag, "// "+e.Notes)
	} else {
		p.Add(e.Name, e.Type, tag)
	}

	return p.Generates()[0]
}

// SetCreatTableStr Set up SQL create statement, backup use setup create statement, backup use.设置创建语句，备份使用
func (s *GenStruct) SetCreatTableStr(sql string) {
	s.SQLBuildStr = sql
}

// SetTableName Setting the name of struct.设置struct名字
func (s *GenStruct) SetTableName(name string) {
	s.TableName = name
}

// SetStructName Setting the name of struct.设置struct名字
func (s *GenStruct) SetStructName(name string) {
	s.Name = name
}

// SetNotes set the notes.设置注释
func (s *GenStruct) SetNotes(notes string) {
	if len(notes) == 0 {
		notes = "[...]" // default of struct notes(for export ).struct 默认注释(为了导出注释)
	}

	notes = s.Name + " " + notes

	a := strings.Split(notes, "\n")
	var text []string

	for _, v := range a {
		text = append(text, "// "+v)
	}

	s.Notes = strings.Join(text, ";")
}

// AddElement Add one or more elements.添加一个/或多个元素
func (s *GenStruct) AddElement(e ...GenElement) {
	s.Em = append(s.Em, e...)
}

// Generates Get the result data.获取结果数据
func (s *GenStruct) Generates() []string {
	var p generate.PrintAtom
	p.Add(s.Notes)
	p.Add("type", s.Name, "struct {")
	mp := make(map[string]bool, len(s.Em))
	for _, v := range s.Em {
		if !mp[v.Name] {
			mp[v.Name] = true
			p.Add(v.Generate())
		}
	}
	p.Add("}")

	return p.Generates()
}

// SetPackage Defining package names.定义包名
func (p *GenPackage) SetPackage(name string) {
	p.Name = name
}

// AddImport Add import by type.通过类型添加import
func (p *GenPackage) AddImport(imp string) {
	if p.Imports == nil {
		p.Imports = make(map[string]string)
	}
	p.Imports[imp] = imp
}

// AddStruct Add a structure.添加一个结构体
func (p *GenPackage) AddStruct(st GenStruct) {
	p.Structs = append(p.Structs, st)
}

// Generate Get the result data.获取结果数据
func (p *GenPackage) Generate() string {
	p.genImport()

	var pa generate.PrintAtom
	pa.Add("// Auto generated. DO NOT EDIT IT.")
	pa.Add("// Auto generated. DO NOT EDIT IT.")
	pa.Add("// Auto generated. DO NOT EDIT IT.")
	pa.Add("\n")

	pa.Add("package", p.Name)
	// add import
	if p.Imports != nil {
		pa.Add("import (")
		for _, v := range p.Imports {
			pa.Add(v)
		}
		pa.Add(")")
	}

	for _, v := range p.Structs {
		for _, v1 := range v.Generates() {
			pa.Add(v1)
		}
	}

	// output.输出
	strOut := ""
	for _, v := range pa.Generates() {
		strOut += v + "\n"
	}

	return strOut
}

// compensate and import .获取结果数据
func (p *GenPackage) genImport() {
	for _, v := range p.Structs {
		for _, v1 := range v.Em {
			if v2, ok := cnf.EImportsHead[v1.Type]; ok {
				if len(v2) > 0 {
					p.AddImport(v2)
				}
			}
		}
	}
}
