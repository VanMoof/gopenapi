package interpret_test

import (
	"github.com/VanMoof/gopenapi/interpret"
	"github.com/VanMoof/gopenapi/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestASTInterpreter_Info(t *testing.T) {
	a := assert.New(t)

	file, openError := os.Open("./_test_files/main_with_info.go")
	a.NoError(openError)

	root := models.Root{}
	interpreter := &interpret.ASTInterpreter{}
	a.NoError(interpreter.InterpretFile(file, &root))

	a.Equal("1.0", root.Info.Version)
	a.Equal("The App Name", root.Info.Title)
	a.Equal("The app description", root.Info.Description)
	a.Equal("Jimbob Jones", root.Info.Contact.Name)
	a.Equal("https://jones.com", root.Info.Contact.URL)
	a.Equal("jimbob@jones.com", root.Info.Contact.Email)
	a.Equal("Apache 2.0", root.Info.License.Name)
	a.Equal("https://www.apache.org/licenses/LICENSE-2.0.html", root.Info.License.URL)
}

func TestASTInterpreter_Path(t *testing.T) {
	a := assert.New(t)

	file, openError := os.Open("./_test_files/func_with_path.go")
	a.NoError(openError)

	root := models.Root{}
	interpreter := &interpret.ASTInterpreter{}
	a.NoError(interpreter.InterpretFile(file, &root))

	a.Equal("The default response of \"ping\"", root.Paths["/ping"].Get.Responses["200"].Description)
	a.Equal("pong", root.Paths["/ping"].Get.Responses["200"].Content["text/plain"].Example)
}
func TestASTInterpreter_Models(t *testing.T) {
	a := assert.New(t)

	file, openError := os.Open("./_test_files/structs_with_models.go")
	a.NoError(openError)

	root := models.Root{}
	interpreter := &interpret.ASTInterpreter{}
	a.NoError(interpreter.InterpretFile(file, &root))
	schemas := root.Components.Schemas

	rootModel := schemas["rootModel"]
	a.Equal("object", rootModel.Type)
	a.Equal("integer", rootModel.Properties["intField"].Type)
	a.Equal("int64", rootModel.Properties["intField"].Format)
	a.Equal("string", rootModel.Properties["stringField"].Type)
	a.Equal("array", rootModel.Properties["subModels"].Type)
	a.Equal("#/components/schemas/subModel", rootModel.Properties["subModels"].Items.Ref)

	subModel := schemas["subModel"]
	a.Equal("object", subModel.Type)
	a.Equal("number", subModel.Properties["floatField"].Type)
	a.Equal("double", subModel.Properties["floatField"].Format)
	a.Equal("object", subModel.Properties["subSubModel"].Type)
	a.Equal("#/components/schemas/subSubModel", subModel.Properties["subSubModel"].AdditionalProperties.(*models.Schema).Ref)

	subSubModel := schemas["subSubModel"]
	a.Equal("object", subSubModel.Type)
	a.Equal("boolean", subSubModel.Properties["boolField"].Type)
	a.Equal("#/components/schemas/aliasedSubs", subSubModel.Properties["aliased"].Ref)

	aliasedSubs := schemas["aliasedSubs"]
	a.Equal("array", aliasedSubs.Type)
	a.Equal("#/components/schemas/aliasedSub", aliasedSubs.Items.Ref)

	aliasedSub := schemas["aliasedSub"]
	a.Equal("object", aliasedSub.Type)
	a.Equal("string", aliasedSub.Properties["timeField"].Type)
	a.Equal("date-time", aliasedSub.Properties["timeField"].Format)
}
