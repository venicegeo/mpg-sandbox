package dsl

import "fmt"

//---------------------------------------------------------------------

type TypeNode interface {
	String() string
}

//---------------------------------------------------------------------------

type TypeNodeName struct {
	Name string
}

func NewTypeNodeName(name string) *TypeNodeName {
	n := &TypeNodeName{
		Name: name,
	}
	return n
}

func (n *TypeNodeName) String() string {
	return fmt.Sprintf("%s", n.Name)
}

//---------------------------------------------------------------------------

type TypeNodeMap struct {
	KeyType   TypeNode
	ValueType TypeNode
}

func NewTypeNodeMap(keyType TypeNode, valueType TypeNode) *TypeNodeMap {
	n := &TypeNodeMap{
		KeyType:   keyType,
		ValueType: valueType,
	}
	return n
}

func (t *TypeNodeMap) String() string {
	return fmt.Sprintf("MAP[%v]%v", t.KeyType, t.ValueType)
}

//---------------------------------------------------------------------------

type TypeNodeStruct struct {
	Fields map[string]TypeNode
}

func NewTypeNodeStruct() *TypeNodeStruct {
	n := &TypeNodeStruct{
		Fields: map[string]TypeNode{},
	}
	return n
}

func (t *TypeNodeStruct) String() string {
	s := ""
	for k := range t.Fields {
		if s != "" {
			s += ", "
		}
		s += fmt.Sprintf("%v", k)
	}
	return fmt.Sprintf("STRUCT(%s)", s)
}

//---------------------------------------------------------------------------

type TypeNodeField struct {
	Name string
	Node TypeNode
}

func NewTypeNodeField(name string, node TypeNode) *TypeNodeField {
	n := &TypeNodeField{
		Node: node,
		Name: name,
	}
	return n
}

func (t *TypeNodeField) String() string {
	return fmt.Sprintf("FIELD(%s,%v)", t.Name, t.Node)
}

//---------------------------------------------------------------------------

type TypeNodeArray struct {
	ElemType TypeNode
	Size     int
}

func NewTypeNodeArray(elemType TypeNode, size int) *TypeNodeArray {
	n := &TypeNodeArray{
		ElemType: elemType,
		Size:     size,
	}
	return n
}

func (t *TypeNodeArray) String() string {
	return fmt.Sprintf("ARRAY(%d, %v)", t.Size, t.ElemType)
}

//---------------------------------------------------------------------------

type TypeNodeSlice struct {
	ElemType TypeNode
}

func NewTypeNodeSlice(elemType TypeNode) *TypeNodeSlice {
	n := &TypeNodeSlice{
		ElemType: elemType,
	}
	return n
}

func (t *TypeNodeSlice) String() string {
	return fmt.Sprintf("SLICE(%v)", t.ElemType)
}

//---------------------------------------------------------------------------

type TypeNodeInt struct {
}

func NewTypeNodeInt() *TypeNodeInt {
	n := &TypeNodeInt{}
	return n
}

func (t *TypeNodeInt) String() string {
	return fmt.Sprintf("INT")
}

//---------------------------------------------------------------------------

type TypeNodeFloat struct {
}

func NewTypeNodeFloat() *TypeNodeFloat {
	n := &TypeNodeFloat{}
	return n
}

func (t *TypeNodeFloat) String() string {
	return fmt.Sprintf("FLOAT")
}

//---------------------------------------------------------------------------

type TypeNodeBool struct {
}

func NewTypeNodeBool() *TypeNodeBool {
	return &TypeNodeBool{}
}

func (t *TypeNodeBool) String() string {
	return "BOOL"
}

//---------------------------------------------------------------------------

type TypeNodeString struct {
}

func NewTypeNodeString() *TypeNodeString {
	return &TypeNodeString{}
}

func (t *TypeNodeString) String() string {
	return "STRING"
}

//---------------------------------------------------------------------------