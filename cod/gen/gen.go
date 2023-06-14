package main

import (
	"fmt"
	"bytes"
	"os"
	"strings"
	// "math"
	"io/fs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/format"
	"text/template"
	"strconv"


	// _ "embed"
)

func main() {
	generatePackage(".")
}

func generatePackage(dir string) {
	// We essentially parse the directory into the fileset and list of packages
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Then we loop over the packages
	for _, pkg := range packages {
		fmt.Println("Parsing", pkg.Name)
		// tokenStart := pkg.Pos()

		// We build a blog visitor (which implements the ast.Visitor interface).
		// This will be used to walk the entire AST!
		bv := &BlogVisitor{
			buf: &bytes.Buffer{},
			pkg: pkg,
			fset: fset,
			// lastCommentPos: &tokenStart,
			structs: make(map[string]StructData),
		}

		// We start walking our BlogVisitor `bv` through the AST in a depth-first way.
		ast.Walk(bv, pkg)

		bv.Output("cod_encode.go")
	}
}

// func formatFunc(buf *bytes.Buffer, fset *token.FileSet, decl ast.FuncDecl, cGroups []*ast.CommentGroup) {
// 	decl.Doc = nil // nil the Doc field so that we don't print it

// 	// Build a CommentedNode. This is important, if you don't attach the comment
// 	// group to the node then the comments inside the function will be removed!
// 	commentedNode := printer.CommentedNode{
// 		Node:     &decl,
// 		Comments: cGroups,
// 	}
// 	formatNode(buf, fset, &commentedNode)
// }

func formatGen(buf *bytes.Buffer, fset *token.FileSet, decl ast.GenDecl, cGroups []*ast.CommentGroup) (StructData, bool) {
	structData := StructData{}

	// commentedNode := printer.CommentedNode{
	// 	Node:     &decl,
	// 	Comments: cGroups,
	// }

	// Skip everything that isn't a type. we can only generate for types
	if decl.Tok != token.TYPE {
		decl.Doc = nil // nil the Doc field so that we don't print it
		return structData, false
	}

	// if decl.Tok == token.IMPORT || decl.Tok == token.TYPE{
	// 	decl.Doc = nil // nil the Doc field so that we don't print it
	// } else if decl.Tok == token.CONST || decl.Tok == token.VAR {
	// 	// Don't nil the documentation
	// }

	// formatNode(buf, fset, &commentedNode)


	fields := make([]Field, 0)

	for _, spec := range decl.Specs {

		switch s := spec.(type) {
		case *ast.TypeSpec:
			directive, directiveCSV := getDirective(decl)
			if directive == DirectiveNone {
				return structData, false
			}

			fmt.Println("TypeSpec: ", s.Name.Name)
			structData.Name = s.Name.Name
			structData.Directive = directive
			structData.DirectiveCSV = directiveCSV

			fmt.Printf("Struct Type: %T\n", s.Type)
			sType, ok := s.Type.(*ast.StructType)
			if !ok {
				// Not a struct, then its an alias. So handle that if we can
				name := "t"
				idxDepth := 0
				field := generateField(name, idxDepth+1, s.Type)
				fields = append(fields, &AliasField{
					Name: name,
					AliasType: s.Name.Name,
					Field: field,
					IndexDepth: idxDepth,
				})
				continue
			}

			for _, f := range sType.Fields.List {
				for _, n := range f.Names {
					fmt.Println("Field: ", n.Name, f.Type, f.Tag)
					fmt.Printf("%T\n", f.Type)

					field := generateField("t." + n.Name, 0, f.Type)
					if f.Tag != nil {
						field.SetTag(f.Tag.Value)
					}

					fields = append(fields, field)
				}
			}
		}
	}
	structData.Fields = fields

	return structData, true
}

func generateField(name string, idxDepth int, node ast.Node) Field {
	switch expr := node.(type) {
	case *ast.Ident:
		fmt.Println("Ident: ", expr.Name)

		if expr.Obj != nil {
			// Then type alias?
			fmt.Println("AAAAAAA", name, expr.Obj.Kind, expr.Obj.Name)
			fmt.Printf("Obj: %T -- %T\n", expr.Obj, expr.Obj.Decl)
		}
		field := &BasicField{
			Name: name,
			Type: expr.Name,
		}

		return field

	// case *ast.StarExpr:
	// 	// fmt.Println("StarExpr: ", expr.Name)
	// 	fmt.Printf("STAR %T\n", expr.X)
	// 	field := generateField(name, 0, expr.X) // TODO: idxDepth???
	// 	return PointerField{
	// 		Field: field,
	// 	}

	case *ast.ArrayType:
		fmt.Printf("ARRAY %T %T\n", expr.Len, expr.Elt)

		if expr.Len == nil {
			idxString := fmt.Sprintf("[i%d]", idxDepth)
			field := generateField(name + idxString, idxDepth + 1, expr.Elt)
			return &SliceField{
				Name: name,
				// Type: field.GetType(),
				Field: field,
				IndexDepth: idxDepth,
			}
		} else {
			idxString := fmt.Sprintf("[i%d]", idxDepth)
			field := generateField(name + idxString, idxDepth + 1, expr.Elt)

			lString := expr.Len.(*ast.BasicLit).Value
			length, err := strconv.Atoi(lString)
			if err != nil { panic("ERR") }
			return &ArrayField{
				Name: name,
				Len: length,
				Field: field,
				IndexDepth: idxDepth,
			}
		}

	case *ast.MapType:
		fmt.Printf("MAP %T %T\n", expr.Key, expr.Value)
		keyString := fmt.Sprintf("[k%d]", idxDepth)
		valString := fmt.Sprintf("[v%d]", idxDepth)
		key := generateField(name + keyString, idxDepth + 1, expr.Key)
		val := generateField(name + valString, idxDepth + 1, expr.Value)
		return &MapField{
			Name: name,
			Key: key,
			Val: val,
			IndexDepth: idxDepth,
		}
	case *ast.SelectorExpr:
		// Note: anything that is a selector expression (ie phy.Position) is guaranteed to be a struct. so it must implement the required struct interface
		fmt.Printf("SEL: %T\n", expr.X)
		field := &BasicField{
			Name: name,
			Type: "UNKNOWN_SELECTOR_EXPR", // This will force it to resolve to the struct marshaller
		}
		return field

	default:
		panic(fmt.Sprintf("unknown type %T", expr))
	}

	return nil
}

// // just make this function take in the field, then have it return the formatted BasicField or PointerField or MapField or ArrayField (or whatever as a response)
// func getFieldType(t ast.Expr) (string) {
// 	switch expr := t.(type) {
// 	case *ast.Ident:
// 		fmt.Println("Ident: ", expr.Name)
// 		// if expr.Name == "string" {
// 		// 	return FieldString
// 		// } else if expr.Name == "uint8" {
// 		// 	return FieldUint8
// 		// }
// 		return expr.Name
// 	case *ast.StarExpr:
// 		// fmt.Println("StarExpr: ", expr.Name)
// 		fmt.Printf("STAR %T\n", expr.X)
// 		return getFieldType(expr.X)
// 	case *ast.ArrayType:
// 		fmt.Printf("ARRAY %T %T\n", expr.Len, expr.Elt)
// 		return getFieldType(expr.Elt)

// 	case *ast.MapType:
// 		fmt.Printf("MAP %T %T\n", expr.Key, expr.Value)
// 	default:
// 		panic(fmt.Sprintf("unknown type %T", expr))
// 	}
// 	return "none"
// }

func getDirective(t ast.GenDecl) (DirectiveType, []string) {
	if t.Doc == nil {
		return DirectiveNone, nil
	}

	for _, c := range t.Doc.List {
		fmt.Println("Doc: ", c.Text)
		if strings.HasPrefix(c.Text, "//cod:struct") {
			after, found := strings.CutPrefix(c.Text, "//cod:struct")
			if !found { panic("BUG") }
			unionNames := strings.Split(after, ",")
			for i := range unionNames {
				unionNames[i] = strings.TrimSpace(unionNames[i])
			}

			return DirectiveStruct, unionNames
		} else if strings.HasPrefix(c.Text, "//cod:union") {
			after, found := strings.CutPrefix(c.Text, "//cod:union")
			if !found { panic("BUG") }
			unionNames := strings.Split(after, ",")
			for i := range unionNames {
				unionNames[i] = strings.TrimSpace(unionNames[i])
			}

			return DirectiveUnion, unionNames
		}
	}
	return DirectiveNone, nil
}

func formatNode(buf *bytes.Buffer, fset *token.FileSet, node any) {
	config := printer.Config{
		Mode:     printer.UseSpaces,
		Tabwidth: 2,
	}

	err := config.Fprint(buf, fset, node)

	if err != nil {
		panic(err)
	}
}

type BlogVisitor struct {
	buf *bytes.Buffer // The buffered blog output
	pkg *ast.Package  // The package that we are processing
	fset *token.FileSet // The fileset of the package we are processing
	file *ast.File // The file we are currently processing (Can be nil if we haven't started processing a file yet!)
	cmap ast.CommentMap // The comment map of the file we are processing

	// lastCommentPos *token.Pos // The token.Pos of the last comment or node that we processed
	structs map[string]StructData
}

func (v *BlogVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil { return nil }

	// If we are a package, then just keep searching
	_, ok := node.(*ast.Package)
	if ok { return v }

	// If we are a file, then store some data in the visitor so we can use it later
	file, ok := node.(*ast.File)
	if ok {
		v.file = file
		v.cmap = ast.NewCommentMap(v.fset, file, file.Comments)
		return v
	}

	// If we are a function, do the function formatting
	_, ok = node.(*ast.FuncDecl)
	if ok {
		return nil // Skip: we don't handle funcs
	}

	gen, ok := node.(*ast.GenDecl)
	if ok {
		cgroups := v.cmap.Filter(gen).Comments()
		sd, ok := formatGen(v.buf, v.fset, *gen, cgroups)
		if ok {
			v.structs[sd.Name] = sd
		}


		return nil
	}

	// If all else fails, then keep looking
	return v
}

type StructData struct {
	Name string
	Directive DirectiveType
	DirectiveCSV []string
	Fields []Field
}

func (s *StructData) WriteStructMarshal(buf *bytes.Buffer) {
	if s.Directive == DirectiveStruct {
		for _, f := range s.Fields {
			f.WriteMarshal(buf)
		}
	} else if s.Directive == DirectiveUnion {
		innerBuf := new(bytes.Buffer)
		for tag, name := range s.DirectiveCSV {
			err := BasicTemp.ExecuteTemplate(innerBuf, "union_case_marshal", map[string]any{
				"Type": name,
				"Tag": tag+1,
				// "InnerCode": string(innerInnerBuf.Bytes()),
				// For now I'm just going to use the requirement that you can only add items to the union that implement the EncodeCod function pair. But you could fix this and let in primitives. It just gets hard and isnt as useful for me right now
				// "InnerCode": "t.EncodeCod(buf)",
			})
			if err != nil { panic(err) }
		}

		err := BasicTemp.ExecuteTemplate(buf, "union_marshal", map[string]any{
			"InnerCode": string(innerBuf.Bytes()),
		})
		if err != nil { panic(err) }

	}
}

func (s *StructData) WriteStructUnmarshal(buf *bytes.Buffer) {
	if s.Directive == DirectiveStruct {
		for _, f := range s.Fields {
			f.WriteUnmarshal(buf)
		}
	} else if s.Directive == DirectiveUnion {
		innerBuf := new(bytes.Buffer)
		for tag, name := range s.DirectiveCSV {
			err := BasicTemp.ExecuteTemplate(innerBuf, "union_case_unmarshal", map[string]any{
				"Type": name,
				"Tag": tag+1,
				// "InnerCode": "t.EncodeCod(buf)",
			})
			if err != nil { panic(err) }
		}

		err := BasicTemp.ExecuteTemplate(buf, "union_unmarshal", map[string]any{
			"InnerCode": string(innerBuf.Bytes()),
		})
		if err != nil { panic(err) }
	}
}

type Field interface {
	WriteMarshal(*bytes.Buffer)
	WriteUnmarshal(*bytes.Buffer)
	SetTag(string)
	SetName(string)
	GetType() string
}

type BasicField struct {
	Name string
	Type string
	Tag string
}

func (f *BasicField) SetName(name string) {
	f.Name = name
}

func (f *BasicField) SetTag(tag string) {
	f.Tag = tag
}
func (f *BasicField) GetType() string {
	return f.Type
}

func (f BasicField) WriteMarshal(buf *bytes.Buffer) {
	templateName := fmt.Sprintf("%s_%s_marshal", "reg", f.Type)
	err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
		"Name": f.Name,
	})

	if err != nil {
		fmt.Println("Couldn't find type, assuming its a struct: ", f.Name)
		templateName := fmt.Sprintf("%s_%s_marshal", "reg", "struct")
		err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
			"Name": f.Name,
		})
		if err != nil { panic(err) }
	}

	// pointerStar := "reg"
	// if f.Pointer { pointerStar = "ptr" }

	// templateName := fmt.Sprintf("%s_%s_marshal", pointerStar, f.Type)
	// err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
	// 	"Name": f.Name,
	// })
	// if err != nil {
	// 	fmt.Println("Couldn't find type, assuming its a struct: ", f.Name)
	// 	templateName := fmt.Sprintf("%s_%s_marshal", pointerStar, "struct")
	// 	err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
	// 		"Name": f.Name,
	// 	})
	// 	if err != nil { panic(err) }
	// }
}

func (f BasicField) WriteUnmarshal(buf *bytes.Buffer) {
	templateName := fmt.Sprintf("%s_%s_unmarshal", "reg", f.Type)
	err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
		"Name": f.Name,
	})
	if err != nil {
		fmt.Println("Couldn't find type, assuming its a struct: ", f.Name)
		templateName := fmt.Sprintf("%s_%s_unmarshal", "reg", "struct")
		err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
			"Name": f.Name,
		})
		if err != nil { panic(err) }
	}


	// pointerStar := "reg"
	// if f.Pointer { pointerStar = "ptr" }

	// templateName := fmt.Sprintf("%s_%s_unmarshal", pointerStar, f.Type)
	// err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
	// 	"Name": f.Name,
	// })
	// if err != nil {
	// 	fmt.Println("Couldn't find type, assuming its a struct: ", f.Name)
	// 	templateName := fmt.Sprintf("%s_%s_unmarshal", pointerStar, "struct")
	// 	err := BasicTemp.ExecuteTemplate(buf, templateName, map[string]any{
	// 		"Name": f.Name,
	// 	})
	// 	if err != nil { panic(err) }
	// }
}

type ArrayField struct {
	Name string
	Field Field
	Len int
	Tag string
	IndexDepth int
}

func (f *ArrayField) SetName(name string) {
	f.Name = name
}
func (f *ArrayField) SetTag(tag string) {
	f.Tag = tag
}
func (f *ArrayField) GetType() string {
	return fmt.Sprintf("[%d]%s", f.Len, f.Field.GetType())
}

func (f ArrayField) WriteMarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	f.Field.WriteMarshal(innerBuf)


	err := BasicTemp.ExecuteTemplate(buf, "array_marshal", map[string]any{
		"Name": f.Name,
		"Index": fmt.Sprintf("i%d", f.IndexDepth),
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }
}


func (f ArrayField) WriteUnmarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	f.Field.WriteUnmarshal(innerBuf)

	err := BasicTemp.ExecuteTemplate(buf, "array_unmarshal", map[string]any{
		"Name": f.Name,
		"Index": fmt.Sprintf("i%d", f.IndexDepth),
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }

}

type SliceField struct {
	Name string
	// Type string
	Field Field
	Tag string
	IndexDepth int
}

func (f *SliceField) SetName(name string) {
	f.Name = name
}
func (f *SliceField) SetTag(tag string) {
	f.Tag = tag
}
func (f *SliceField) GetType() string {
	return fmt.Sprintf("[]%s", f.Field.GetType())
}

func (f SliceField) WriteMarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	idxVar := fmt.Sprintf("i%d", f.IndexDepth)
	f.Field.SetName(fmt.Sprintf("%s[%s]", f.Name, idxVar))
	f.Field.WriteMarshal(innerBuf)

	err := BasicTemp.ExecuteTemplate(buf, "slice_marshal", map[string]any{
		"Name": f.Name,
		"Type": f.Field.GetType(),
		"Index": idxVar,
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }
}


func (f SliceField) WriteUnmarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	varName := fmt.Sprintf("value%d", f.IndexDepth)
	f.Field.SetName(varName)
	f.Field.WriteUnmarshal(innerBuf)

	fmt.Println("GETTYPE: ", f.Field.GetType())
	err := BasicTemp.ExecuteTemplate(buf, "slice_unmarshal", map[string]any{
		"Name": f.Name,
		"VarName": varName,
		"Type": f.Field.GetType(),
		"Index": fmt.Sprintf("i%d", f.IndexDepth),
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }

}

type MapField struct {
	Name string
	Key Field
	Val Field
	Tag string
	IndexDepth int
}

func (f *MapField) SetName(name string) {
	f.Name = name
}
func (f *MapField) SetTag(tag string) {
	f.Tag = tag
}
func (f *MapField) GetType() string {
	return fmt.Sprintf("map[%s]%s", f.Key.GetType(), f.Val.GetType())
}

func (f MapField) WriteMarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)

	keyIdxName := fmt.Sprintf("k%d", f.IndexDepth)
	f.Key.SetName(keyIdxName)
	f.Key.WriteMarshal(innerBuf)

	valIdxName := fmt.Sprintf("v%d", f.IndexDepth)
	f.Val.SetName(valIdxName)
	f.Val.WriteMarshal(innerBuf)

	err := BasicTemp.ExecuteTemplate(buf, "map_marshal", map[string]any{
		"Name": f.Name,
		"KeyIdx": keyIdxName,
		"ValIdx": valIdxName,
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }
}


func (f MapField) WriteUnmarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	keyVarName := fmt.Sprintf("key%d", f.IndexDepth)
	f.Key.SetName(keyVarName)
	f.Key.WriteUnmarshal(innerBuf)

	valVarName := fmt.Sprintf("val%d", f.IndexDepth)
	f.Val.SetName(valVarName)
	f.Val.WriteUnmarshal(innerBuf)

	fmt.Println("GETTYPE: ", f.GetType(), f.Key.GetType(), f.Val.GetType())
	err := BasicTemp.ExecuteTemplate(buf, "map_unmarshal", map[string]any{
		"Name": f.Name,
		"Type": f.GetType(),
		"KeyVar": keyVarName,
		"ValVar": valVarName,
		"KeyType": f.Key.GetType(),
		"ValType": f.Val.GetType(),
		"InnerCode": string(innerBuf.Bytes()),

		"Index": fmt.Sprintf("i%d", f.IndexDepth),
	})
	if err != nil { panic(err) }

}

type AliasField struct {
	Name string
	AliasType string
	Field Field
	Tag string
	IndexDepth int
}

func (f *AliasField) SetName(name string) {
	f.Name = name
}
func (f *AliasField) SetTag(tag string) {
	f.Tag = tag
}
func (f *AliasField) GetType() string {
	return fmt.Sprintf("%s", f.Field.GetType())
}

func (f AliasField) WriteMarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)

	valName := fmt.Sprintf("value%d", f.IndexDepth)
	f.Field.SetName(valName)
	f.Field.WriteMarshal(innerBuf)

	err := BasicTemp.ExecuteTemplate(buf, "alias_marshal", map[string]any{
		"Name": f.Name,
		"AliasType": f.Name,
		"Type": f.GetType(),
		"ValName": valName,
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }
}


func (f AliasField) WriteUnmarshal(buf *bytes.Buffer) {
	innerBuf := new(bytes.Buffer)
	valName := fmt.Sprintf("value%d", f.IndexDepth)
	f.Field.SetName(valName)
	f.Field.WriteUnmarshal(innerBuf)

	fmt.Println("ALIAS_GETTYPE: ", f.GetType(), f.Field.GetType())
	err := BasicTemp.ExecuteTemplate(buf, "alias_unmarshal", map[string]any{
		"Name": f.Name,
		"AliasType": f.AliasType,
		"Type": f.GetType(),
		"ValName": valName,
		"ValType": f.Field.GetType(),
		"InnerCode": string(innerBuf.Bytes()),
	})
	if err != nil { panic(err) }

}

// type FieldType uint16
// const (
// 	FieldNone FieldType = iota
// 	FieldUint8
// 	FieldUint16
// 	FieldUint32
// 	FieldUint64
// 	FieldInt8
// 	FieldInt16
// 	FieldInt32
// 	FieldInt64

// 	FieldString
// 	FieldStruct
// )
type DirectiveType uint8
const (
	DirectiveNone DirectiveType = iota
	DirectiveStruct
	DirectiveUnion
)


func (v *BlogVisitor) Output(filename string) {
	marshal, err := template.New("marshal").Parse(`
func (t {{.Name}})EncodeCod(buf *cod.Buffer) {
{{.MarshalCode}}
}
`)
	if err != nil { panic(err) }

	unmarshal, err := template.New("unmarshal").Parse(`
func (t *{{.Name}})DecodeCod(buf *cod.Buffer) error {
var err error

{{.MarshalCode}}

return err
}
`)
	if err != nil { panic(err) }

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("package " + v.pkg.Name)
	buf.WriteString(`
import (
	"github.com/unitoftime/flow/cod"
)`)


	marshBuf := bytes.NewBuffer([]byte{})
	unmarshBuf := bytes.NewBuffer([]byte{})
	for _, sd := range v.structs {
		marshBuf.Reset()
		unmarshBuf.Reset()

		fmt.Println("Struct: ", sd)

		// Write the marshal code
		sd.WriteStructMarshal(marshBuf)

		// Write the unmarshal code
		sd.WriteStructUnmarshal(unmarshBuf)

		// Write the encode func
		err = marshal.Execute(buf, map[string]any{
			"Name": sd.Name,
			"MarshalCode": string(marshBuf.Bytes()),
		})
		if err != nil { panic(err) }

		// Write the decode func
		err = unmarshal.Execute(buf, map[string]any{
			"Name": sd.Name,
			"MarshalCode": string(unmarshBuf.Bytes()),
		})
		if err != nil { panic(err) }
	}

	for _, sd := range v.structs {
		if sd.Directive != DirectiveUnion { continue }

		// Create constructors, getters, setters per union type
		err := BasicTemp.ExecuteTemplate(buf, "union_getter", map[string]any{
			"Name": sd.Name,
		})
		if err != nil { panic(err) }
		err = BasicTemp.ExecuteTemplate(buf, "union_setter", map[string]any{
			"Name": sd.Name,
		})
		if err != nil { panic(err) }
		err = BasicTemp.ExecuteTemplate(buf, "union_constructor", map[string]any{
			"Name": sd.Name,
		})
		if err != nil { panic(err) }
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		err = os.WriteFile(filename, buf.Bytes(), fs.ModePerm)
		if err != nil {
			panic(err)
		}
		return
	}

	err = os.WriteFile(filename, formatted, fs.ModePerm)
	// err = os.WriteFile(filename, buf.Bytes(), fs.ModePerm)
	if err != nil {
		panic(err)
	}
}


