package typedef

import "reflect"

type Field struct {
	// 中文名称
	NameCN string
	// 备注
	Remark string
	// 渲染json用的名称
	JsonName string
	// select SQL 后 row 内的 key
	ColName string
	// Go语言内部的类型
	Kind reflect.Kind
}
