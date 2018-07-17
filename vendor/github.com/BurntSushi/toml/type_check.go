package toml





type tomlType interface {
	typeString() string
}


func typeEqual(t1, t2 tomlType) bool {
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.typeString() == t2.typeString()
}

func typeIsHash(t tomlType) bool {
	return typeEqual(t, tomlHash) || typeEqual(t, tomlArrayHash)
}

type tomlBaseType string

func (btype tomlBaseType) typeString() string {
	return string(btype)
}

func (btype tomlBaseType) String() string {
	return btype.typeString()
}

var (
	tomlInteger   tomlBaseType = "Integer"
	tomlFloat     tomlBaseType = "Float"
	tomlDatetime  tomlBaseType = "Datetime"
	tomlString    tomlBaseType = "String"
	tomlBool      tomlBaseType = "Bool"
	tomlArray     tomlBaseType = "Array"
	tomlHash      tomlBaseType = "Hash"
	tomlArrayHash tomlBaseType = "ArrayHash"
)






func (p *parser) typeOfPrimitive(lexItem item) tomlType {
	switch lexItem.typ {
	case itemInteger:
		return tomlInteger
	case itemFloat:
		return tomlFloat
	case itemDatetime:
		return tomlDatetime
	case itemString:
		return tomlString
	case itemMultilineString:
		return tomlString
	case itemRawString:
		return tomlString
	case itemRawMultilineString:
		return tomlString
	case itemBool:
		return tomlBool
	}
	p.bug("Cannot infer primitive type of lex item '%s'.", lexItem)
	panic("unreachable")
}






func (p *parser) typeOfArray(types []tomlType) tomlType {
	
	if len(types) == 0 {
		return tomlArray
	}

	theType := types[0]
	for _, t := range types[1:] {
		if !typeEqual(theType, t) {
			p.panicf("Array contains values of type '%s' and '%s', but "+
				"arrays must be homogeneous.", theType, t)
		}
	}
	return tomlArray
}
