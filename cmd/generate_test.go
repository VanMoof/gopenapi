package cmd_test

import (
	"bytes"
	"encoding/json"
	"github.com/VanMoof/gopenapi/cmd"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGenerateSpec_YAMLFile(t *testing.T) {
	a := assert.New(t)

	tempFile, tempFileError := ioutil.TempFile("", "*.yaml")
	a.NoError(tempFileError)
	a.NoError(cmd.GenerateSpec("yaml", tempFile.Name(), []string{"../interpret/_test_files"}))

	decoded := map[string]interface{}{}
	a.NoError(yaml.NewDecoder(tempFile).Decode(&decoded))
	a.Equal("3.0.2", decoded["openapi"])
}

func TestGenerateSpec_JSONFile(t *testing.T) {
	a := assert.New(t)

	tempFile, tempFileError := ioutil.TempFile("", "*.yaml")
	a.NoError(tempFileError)
	a.NoError(cmd.GenerateSpec("json", tempFile.Name(), []string{"../interpret/_test_files"}))

	decoded := map[string]interface{}{}
	a.NoError(json.NewDecoder(tempFile).Decode(&decoded))
	a.Equal("3.0.2", decoded["openapi"])
}

func TestGenerateSpec_JSONStdout(t *testing.T) {
	a := assert.New(t)

	writeFunc := func() {
		a.NoError(cmd.GenerateSpec("json", "-", []string{"../interpret/_test_files"}))
	}
	assertFunc := func(out string) {
		decoded := map[string]interface{}{}
		a.NoError(json.NewDecoder(strings.NewReader(out)).Decode(&decoded))
		a.Equal("3.0.2", decoded["openapi"])
	}
	withPipedStdOut(writeFunc, assertFunc)
}

func TestResolveOutputWriter_Stdout(t *testing.T) {
	a := assert.New(t)

	writeFunc := func() {
		writer, writerError := cmd.ResolveOutputWriter("-")
		a.NoError(writerError)
		writer.Write([]byte{'a'})
	}
	assertFunc := func(out string) {
		a.Equal("a", out)
	}
	withPipedStdOut(writeFunc, assertFunc)
}

func TestResolveOutputWriter_File(t *testing.T) {
	a := assert.New(t)

	tempFile, tempFileError := ioutil.TempFile("", "*.txt")
	a.NoError(tempFileError)
	writer, writerError := cmd.ResolveOutputWriter(tempFile.Name())
	a.NoError(writerError)
	writer.Write([]byte{'a'})

	content, readError := ioutil.ReadFile(tempFile.Name())
	a.NoError(readError)
	a.Equal("a", string(content))
}

func TestResolveOutputSink_YAML(t *testing.T) {
	a := assert.New(t)

	var buff bytes.Buffer
	closableBuff := &ClosableBuff{b: &buff}
	sink := cmd.ResolveOutputSink("yaml", closableBuff)
	a.NoError(sink.Write(map[string]string{"key": "value"}))
	a.Equal("key: value\n", buff.String())
}

func TestResolveOutputSink_JSON(t *testing.T) {
	a := assert.New(t)

	var buff bytes.Buffer
	closableBuff := &ClosableBuff{b: &buff}
	sink := cmd.ResolveOutputSink("json", closableBuff)
	a.NoError(sink.Write(map[string]string{"key": "value"}))
	a.Equal("{\n  \"key\": \"value\"\n}\n", buff.String())
}

func TestNormalizeInputPath_Absolute(t *testing.T) {
	a := assert.New(t)

	normalized, err := cmd.NormalizeInputPath("/a/b/c")
	a.NoError(err)
	a.Equal("/a/b/c", normalized)
}

func TestNormalizeInputPath_Relative(t *testing.T) {
	a := assert.New(t)

	normalized, err := cmd.NormalizeInputPath("jim/bob")
	a.NoError(err)
	a.Contains(normalized, "jim/bob")
}

//Piping stdout taken from https://stackoverflow.com/questions/10473800
func withPipedStdOut(writeFunc func(), assertFunc func(out string)) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	writeFunc()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assertFunc(out)
}

type ClosableBuff struct {
	b *bytes.Buffer
}

func (c *ClosableBuff) Write(p []byte) (n int, err error) {
	return c.b.Write(p)
}

func (c *ClosableBuff) Close() error {
	return nil
}
