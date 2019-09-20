package generate_test

import (
	"errors"
	"github.com/VanMoof/gopenapi/generate"
	"github.com/VanMoof/gopenapi/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestJSONSink(t *testing.T) {
	a := assert.New(t)
	tempJSON, tempFileError := ioutil.TempFile("", "*.json")
	a.NoError(tempFileError)
	sink := generate.JSONSink{W: tempJSON}

	writeError := sink.Write(map[string]string{"key": "value"})
	a.NoError(writeError)

	writtenContent, readError := ioutil.ReadFile(tempJSON.Name())
	a.NoError(readError)
	a.Equal("{\n  \"key\": \"value\"\n}\n", string(writtenContent))
}

func TestYAMLSink(t *testing.T) {
	a := assert.New(t)
	tempYAML, tempFileError := ioutil.TempFile("", "*.yaml")
	a.NoError(tempFileError)
	sink := generate.YAMLSink{W: tempYAML}

	writeError := sink.Write(map[string]string{"key": "value"})
	a.NoError(writeError)

	writtenContent, readError := ioutil.ReadFile(tempYAML.Name())
	a.NoError(readError)
	a.Equal("key: value\n", string(writtenContent))
}

func TestGoFileVisitor(t *testing.T) {
	a := assert.New(t)
	tempDir, tempDirError := ioutil.TempDir("", "some-dir")
	a.NoError(tempDirError)

	tempGoFile, tempGoFileError := ioutil.TempFile(tempDir, "*.go")
	a.NoError(tempGoFileError)
	_, tempTxtFileError := ioutil.TempFile(tempDir, "*.txt")
	a.NoError(tempTxtFileError)

	var visitedFiles []string
	visitor := &generate.GoFileVisitor{BasePath: tempDir}
	visitError := visitor.VisitFiles(func(filePath string, info os.FileInfo, err error) error {
		visitedFiles = append(visitedFiles, filePath)
		return nil
	})
	a.NoError(visitError)

	a.Len(visitedFiles, 1)
	a.Equal(tempGoFile.Name(), visitedFiles[0])
}

func TestGenerate_FailOnVisitor(t *testing.T) {
	a := assert.New(t)

	visitor := &testFileVisitor{fail: true}
	interpreter := &testInterpreter{}
	sink := &testSink{}

	a.Error(generate.Generate(visitor, interpreter, sink))
	a.True(visitor.called)
	a.False(interpreter.called)
	a.False(sink.called)
}

func TestGenerate_FailOnInterpreter(t *testing.T) {
	a := assert.New(t)

	visitor := &testFileVisitor{}
	interpreter := &testInterpreter{fail: true}
	sink := &testSink{}

	a.Error(generate.Generate(visitor, interpreter, sink))
	a.True(visitor.called)
	a.True(interpreter.called)
	a.False(sink.called)
}

func TestGenerate_FailOnSink(t *testing.T) {
	a := assert.New(t)

	visitor := &testFileVisitor{}
	interpreter := &testInterpreter{}
	sink := &testSink{fail: true}

	a.Error(generate.Generate(visitor, interpreter, sink))
	a.True(visitor.called)
	a.True(interpreter.called)
	a.True(sink.called)
}

func TestGenerate(t *testing.T) {
	a := assert.New(t)

	visitor := &testFileVisitor{}
	interpreter := &testInterpreter{}
	sink := &testSink{}

	a.NoError(generate.Generate(visitor, interpreter, sink))
	a.True(visitor.called)
	a.True(interpreter.called)
	a.True(sink.called)
}

type testFileVisitor struct {
	called bool
	fail   bool
}

func (t *testFileVisitor) VisitFiles(f func(filePath string, info os.FileInfo, err error) error) error {
	t.called = true
	if t.fail {
		return errors.New("something happened")
	}
	return f("", nil, nil)
}

type testInterpreter struct {
	called bool
	fail   bool
}

func (t *testInterpreter) InterpretFile(file *os.File, root *models.Root) error {
	t.called = true
	if t.fail {
		return errors.New("something happened")
	}
	return nil
}

type testSink struct {
	called bool
	fail   bool
}

func (t *testSink) Write(interface{}) error {
	t.called = true
	if t.fail {
		return errors.New("something happened")
	}
	return nil
}
