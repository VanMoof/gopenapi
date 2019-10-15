package interpret

import (
	"fmt"
	"github.com/VanMoof/gopenapi/models"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Interpreter interface {
	InterpretFile(file *os.File, root *models.Root) error
}

type ASTInterpreter struct {
}

func (a *ASTInterpreter) InterpretFile(file *os.File, root *models.Root) error {
	fileSet := token.NewFileSet()

	parsedFile, parseError := parser.ParseFile(fileSet, file.Name(), file, parser.ParseComments)
	if parseError != nil {
		return fmt.Errorf("failed to interpret file %s: %w", file.Name(), parseError)
	}
	return interpretFile(parsedFile, root)
}

func interpretFile(parsedFile *ast.File, root *models.Root) error {
	declarations := parsedFile.Decls
	for _, declaration := range declarations {
		switch declaration.(type) {
		case *ast.FuncDecl:
			err := openAPIBlockFromFunctionDeclaration(declaration.(*ast.FuncDecl), root)
			if err != nil {
				return err
			}
		case *ast.GenDecl:
			openAPIBlockFromGenDeclaration(declaration.(*ast.GenDecl), root)
		}
	}
	return nil
}

func openAPIBlockFromFunctionDeclaration(funcDecl *ast.FuncDecl, root *models.Root) error {
	commentGroup := funcDecl.Doc
	cleanedComment := cleanComment(commentGroup.Text())
	err := commentAsOpenAPIBlock(root, cleanedComment)
	if err != nil {
		return fmt.Errorf("failed to resolve comment as OpenAPI element: %w", err)
	}
	return nil
}

func openAPIBlockFromGenDeclaration(genDecl *ast.GenDecl, root *models.Root) {
	switch genDecl.Tok {
	case token.TYPE:
		openAPIBlockFromTypeDeclaration(genDecl, root)
	case token.CONST, token.VAR:
		openAPIBlockFromConstAndVarDeclaration(genDecl, root)
	}
}

func openAPIBlockFromTypeDeclaration(decl *ast.GenDecl, root *models.Root) {
	commentGroup := decl.Doc
	if !strings.Contains(commentGroup.Text(), "gopenapi:objectSchema") {
		return
	}
	for _, spec := range decl.Specs {
		switch spec.(type) {
		case *ast.TypeSpec:
			openAPIBlockFromTypeSpec(spec.(*ast.TypeSpec), root)
		}
	}
}

func openAPIBlockFromConstAndVarDeclaration(decl *ast.GenDecl, root *models.Root) error {
	cleanedComment := cleanComment(decl.Doc.Text())
	if !strings.HasPrefix(cleanedComment, "gopenapi:parameter") {
		return nil
	}
	if root.Components == nil {
		root.Components = &models.Components{}
	}
	if root.Components.Parameters == nil {
		root.Components.Parameters = map[string]*models.Parameter{}
	}
	cleanedComment = strings.TrimPrefix(cleanedComment, "gopenapi:parameter")
	spec := decl.Specs[0]
	valueSpec := spec.(*ast.ValueSpec)
	parameter := models.Parameter{}
	root.Components.Parameters[valueSpec.Names[0].Name] = &parameter
	basicLit := valueSpec.Values[0].(*ast.BasicLit)
	unquoted, unquoteError := strconv.Unquote(basicLit.Value)
	if unquoteError != nil {
		return unquoteError
	}
	parameter.Name = unquoted

	err := yaml.NewDecoder(strings.NewReader(cleanedComment)).Decode(&parameter)
	if err != nil {
		return fmt.Errorf("failed to decode comment:\n%s\nError: %w", cleanedComment, err)
	}
	return nil
}

func openAPIBlockFromTypeSpec(typeSpec *ast.TypeSpec, root *models.Root) {
	if root.Components == nil {
		root.Components = &models.Components{}
	}
	if root.Components.Schemas == nil {
		root.Components.Schemas = map[string]*models.Schema{}
	}

	newSchema := &models.Schema{
		Type:       "object",
		Properties: map[string]*models.Schema{},
	}
	newSchemaName := lower(typeSpec.Name.Name)

	root.Components.Schemas[newSchemaName] = newSchema
	switch typeSpec.Type.(type) {
	case *ast.StructType:
		structType := typeSpec.Type.(*ast.StructType)
		schemaFieldsFromStructType(structType, newSchema)
	case *ast.ArrayType:
		arrayType := typeSpec.Type.(*ast.ArrayType)
		setSchemaType(newSchema, "array")

		starExpr := arrayType.Elt.(*ast.StarExpr)

		newSchema.Items = &models.Schema{}
		setSchemaType(newSchema.Items, starExpr.X.(*ast.Ident).Name)
	}
}

func structFieldName(structField *ast.Field) string {
	if structField.Tag == nil {
		return lower(structField.Names[0].Name)
	}
	structTag := reflect.StructTag(strings.ReplaceAll(structField.Tag.Value, "`", ""))
	fieldName := structTag.Get("json")
	if fieldName == "-" {
		return ""
	}
	return fieldName
}

func schemaFieldsFromStructType(structType *ast.StructType, newSchema *models.Schema) {
	structFields := structType.Fields
	for _, structField := range structFields.List {
		fieldName := structFieldName(structField)
		if fieldName == "" {
			continue
		}
		newSchema.Properties[fieldName] = &models.Schema{}
		structFieldType := structField.Type
		switch structFieldType.(type) {
		case *ast.SelectorExpr:
			selectorExpr := structFieldType.(*ast.SelectorExpr)
			name := fmt.Sprintf("%s.%s", selectorExpr.X.(*ast.Ident).Name, selectorExpr.Sel.Name)
			setSchemaType(newSchema.Properties[fieldName], name)
		case *ast.Ident:
			setSchemaType(newSchema.Properties[fieldName], structFieldType.(*ast.Ident).Name)
		case *ast.ArrayType:
			setSchemaType(newSchema.Properties[fieldName], "array")

			arrayType := structFieldType.(*ast.ArrayType)
			starExpr := arrayType.Elt.(*ast.StarExpr)

			newSchema.Properties[fieldName].Items = &models.Schema{}
			setSchemaType(newSchema.Properties[fieldName].Items, starExpr.X.(*ast.Ident).Name)
		case *ast.MapType:
			setSchemaType(newSchema.Properties[fieldName], "object")
			mapSchema := &models.Schema{}

			mapType := structFieldType.(*ast.MapType)
			setSchemaType(mapSchema, mapType.Value.(*ast.StarExpr).X.(*ast.Ident).Name)
			newSchema.Properties[fieldName].AdditionalProperties = mapSchema
		}
	}
}

func lower(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func commentAsOpenAPIBlock(r *models.Root, comment string) error {
	types := map[string]func(*models.Root) interface{}{
		"gopenapi:info": func(r *models.Root) interface{} {
			r.Info = &models.Info{}
			return r.Info
		},
		"gopenapi:path": func(r *models.Root) interface{} {
			if r.Paths == nil {
				r.Paths = map[string]*models.PathItem{}
			}
			return &r.Paths
		},
	}
	for blockType, modelPointerRetriever := range types {
		if strings.HasPrefix(comment, blockType) {
			comment = strings.TrimPrefix(comment, blockType)
			modelPointer := modelPointerRetriever(r)
			err := yaml.NewDecoder(strings.NewReader(comment)).Decode(modelPointer)
			if err != nil {
				return fmt.Errorf("failed to decode comment:\n%s\nError: %w", comment, err)
			}
			return nil
		}
	}
	return nil
}

func cleanComment(c string) string {
	return strings.ReplaceAll(strings.TrimSpace(c), "\t", "    ")
}

func setSchemaType(schema *models.Schema, typeName string) {
	switch typeName {
	case "bool":
		schema.Type = "boolean"
	case "int64", "int":
		schema.Type = "integer"
		schema.Format = "int64"
	case "int32", "time.Month":
		schema.Type = "integer"
		schema.Format = "int32"
	case "float64", "float":
		schema.Type = "number"
		schema.Format = "double"
	case "float32":
		schema.Type = "number"
		schema.Format = "float"
	case "string":
		schema.Type = "string"
	case "array":
		schema.Type = "array"
	case "object":
		schema.Type = "object"
	case "time.Time":
		schema.Type = "string"
		schema.Format = "date-time"
	default:
		schema.Ref = "#/components/schemas/" + lower(typeName)
	}
	return
}
