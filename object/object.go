package object

import (
	"TLang/ast"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Type string

const (
	INTEGER   = "INTEGER"
	FLOAT     = "FLOAT"
	BOOLEAN   = "BOOLEAN"
	STRING    = "STRING"
	CHARACTER = "CHARACTER"
	VOID      = "VOID"
	RET       = "RET"
	OUT       = "OUT"
	JUMP      = "JUMP"
	ERR       = "ERR"
	FUNCTION  = "FUNCTION"
	UNDERLINE = "UNDERLINE"
	NATIVE    = "NATIVE"
	ARRAY     = "ARRAY"
	REFERENCE = "REFERENCE"
	HASH      = "HASH"
)

type Object interface {
	Type() Type
	Inspect() string
	Copy() Object
}

type Number interface {
	Object
	NumberObj()
}

type Letter interface {
	Object
	LetterObj() string
}

type LikeFunction interface {
	Object
	LikeFunctionObj()
}

type HashAble interface {
	Object
	HashKey() HashKey
}

type AllocRequired interface {
	DoAlloc(Index Object) (*Object, bool)
	DeAlloc(Index Object) bool
}

type HashKey struct {
	Type  Type
	Value interface{}
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() Type      { return INTEGER }
func (i *Integer) Copy() Object    { return i }
func (i *Integer) NumberObj()      {}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: i.Value}
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string { return fmt.Sprintf("%g", f.Value) }
func (f *Float) Type() Type      { return FLOAT }
func (f *Float) Copy() Object    { return f }
func (f *Float) NumberObj()      {}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() Type      { return BOOLEAN }
func (b *Boolean) Copy() Object    { return b }
func (b *Boolean) HashKey() HashKey {
	return HashKey{Type: b.Type(), Value: b.Value}
}

type Void struct{}

func (v *Void) Inspect() string { return "void" }
func (v *Void) Type() Type      { return VOID }
func (v *Void) Copy() Object    { return v }
func (v *Void) HashAble() HashKey {
	return HashKey{Type: v.Type(), Value: 0}
}

type RetValue struct {
	Value Object
}

func (rv *RetValue) Inspect() string { return rv.Value.Inspect() }
func (rv *RetValue) Type() Type      { return RET }
func (rv *RetValue) Copy() Object {
	println("WARNING: COPY RET VALUE")
	return &RetValue{Value: rv.Copy()}
}

type OutValue struct {
	Value Object
}

func (ov *OutValue) Inspect() string { return ov.Value.Inspect() }
func (ov *OutValue) Type() Type      { return OUT }
func (ov *OutValue) Copy() Object {
	println("WARNING: COPY OUT VALUE")
	return &OutValue{Value: ov.Copy()}
}

type Jump struct{}

func (j *Jump) Inspect() string { return "jump" }
func (j *Jump) Type() Type      { return JUMP }
func (j *Jump) Copy() Object {
	println("WARNING: COPY JUMP")
	return j
}

type Err struct {
	Message string
}

func (err *Err) Inspect() string { return "ERROR: " + err.Message }
func (err *Err) Type() Type      { return ERR }
func (err *Err) Copy() Object {
	println("WARNING: COPY ERR")
	return err
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}
func (f *Function) Type() Type       { return FUNCTION }
func (f *Function) Copy() Object     { return f }
func (f *Function) LikeFunctionObj() {}

type UnderLine struct {
	Body *ast.BlockStatement
	Env  *Environment
}

func (u *UnderLine) Inspect() string {
	var out bytes.Buffer

	out.WriteString("_ ")
	out.WriteString(u.Body.String())

	return out.String()
}
func (u *UnderLine) Type() Type       { return UNDERLINE }
func (u *UnderLine) Copy() Object     { return u }
func (u *UnderLine) LikeFunctionObj() {}

type String struct {
	Value string
}

func (s *String) Inspect() string   { return strconv.Quote(s.Value) }
func (s *String) Type() Type        { return STRING }
func (s *String) Copy() Object      { return s }
func (s *String) LetterObj() string { return s.Value }
func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), Value: s.Value}
}

type Character struct {
	Value rune
}

func (c *Character) Inspect() string   { return "'" + string(c.Value) + "'" }
func (c *Character) Type() Type        { return CHARACTER }
func (c *Character) Copy() Object      { return c }
func (c *Character) LetterObj() string { return string(c.Value) }
func (c *Character) HashKey() HashKey {
	return HashKey{Type: c.Type(), Value: c.Value}
}

type Native struct {
	Fn func(env *Environment, args []Object) Object
}

func (n *Native) Inspect() string  { return "func [Native]" }
func (n *Native) Type() Type       { return NATIVE }
func (n *Native) Copy() Object     { return n }
func (n *Native) LikeFunctionObj() {}

type Array struct {
	Elements []Object
	Copyable bool
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Type() Type { return ARRAY }
func (a *Array) Copy() Object {
	if a.Copyable {
		a.Copyable = false
		return a
	}
	var elements []Object
	for _, e := range a.Elements {
		elements = append(elements, e.Copy())
	}
	return &Array{Elements: elements}
}

type Reference struct {
	Value  *Object
	Origin AllocRequired
	Index  Object
	Const  bool
}

func (r *Reference) Inspect() string {
	var out bytes.Buffer

	if r.Const {
		out.WriteString("Const ")
	}
	out.WriteString("Reference: ")
	out.WriteString((*r.Value).Inspect())

	return out.String()
}
func (r *Reference) Type() Type { return REFERENCE }
func (r *Reference) Copy() Object {
	println("WARNING: COPY REFERENCE")
	return (*r.Value).Copy()
}

type HashPair struct {
	Key   Object
	Value *Object
}

type Hash struct {
	Pairs    map[HashKey]HashPair
	Copyable bool
}

func (h *Hash) DoAlloc(Index Object) (*Object, bool) {
	if hashIndex, ok := Index.(HashAble); ok {
		key := hashIndex.HashKey()
		if _, ok := h.Pairs[key]; !ok {
			var obj Object = nil
			h.Pairs[key] = HashPair{
				Key:   hashIndex,
				Value: &obj,
			}
			return &obj, true
		}
	}
	return nil, false
}
func (h *Hash) DeAlloc(Index Object) bool {
	if hashIndex, ok := Index.(HashAble); ok {
		key := hashIndex.HashKey()
		if _, ok := h.Pairs[key]; ok {
			delete(h.Pairs, key)
			return true
		}
	}
	return false
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), (*pair.Value).Inspect()))
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}
func (h *Hash) Type() Type { return HASH }
func (h *Hash) Copy() Object {
	if h.Copyable {
		h.Copyable = false
		return h
	}
	pairs := make(map[HashKey]HashPair)
	for index, pair := range h.Pairs {
		newVal := (*pair.Value).Copy()
		pairs[index] = HashPair{
			Key:   pair.Key,
			Value: &newVal,
		}
	}

	return &Hash{
		Pairs:    pairs,
		Copyable: false,
	}
}
