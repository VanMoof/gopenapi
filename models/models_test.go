package models_test

import (
	"encoding/json"
	"github.com/VanMoof/gopenapi/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestEncodeRoot_YAML(t *testing.T) {
	r := root()
	encoded, err := yaml.Marshal(r)

	a := assert.New(t)
	a.NoError(err)

	r2 := &models.Root{}
	a.NoError(yaml.Unmarshal(encoded, r2))
	validateRoot(a, r2)
}

func TestEncodeRoot_JSON(t *testing.T) {
	r := root()
	encoded, err := json.Marshal(r)

	a := assert.New(t)
	a.NoError(err)

	r2 := &models.Root{}
	a.NoError(json.Unmarshal(encoded, r2))
	validateRoot(a, r2)
}

func TestEncodeInfo_YAML(t *testing.T) {
	i := info()

	encoded, err := yaml.Marshal(i)

	a := assert.New(t)
	a.NoError(err)

	i2 := &models.Info{}
	a.NoError(yaml.Unmarshal(encoded, i2))

	validateInfo(a, i2)
}

func TestEncodeInfo_JSON(t *testing.T) {
	i := info()

	encoded, err := json.Marshal(i)

	a := assert.New(t)
	a.NoError(err)

	i2 := &models.Info{}
	a.NoError(json.Unmarshal(encoded, i2))

	validateInfo(a, i2)
}

func TestEncodeServer_YAML(t *testing.T) {
	s := server()

	encoded, err := yaml.Marshal(s)

	a := assert.New(t)
	a.NoError(err)

	s2 := &models.Server{}
	a.NoError(yaml.Unmarshal(encoded, s2))

	validateServer(a, s2)
}

func TestEncodeServer_JSON(t *testing.T) {
	s := server()

	encoded, err := json.Marshal(s)

	a := assert.New(t)
	a.NoError(err)

	s2 := &models.Server{}
	a.NoError(json.Unmarshal(encoded, s2))

	validateServer(a, s2)
}

func TestComponents_YAML(t *testing.T) {
	c := components()

	encoded, err := yaml.Marshal(c)

	a := assert.New(t)
	a.NoError(err)

	c2 := &models.Components{}
	a.NoError(yaml.Unmarshal(encoded, c2))

	validateComponents(a, c2)
}

func TestComponents_JSON(t *testing.T) {
	c := components()

	encoded, err := json.Marshal(c)

	a := assert.New(t)
	a.NoError(err)

	c2 := &models.Components{}
	a.NoError(json.Unmarshal(encoded, c2))

	validateComponents(a, c2)
}

func TestMergePaths_PartialDuplicates(t *testing.T) {
	item1 := `
/the/path:
  $ref: ref
  summary: summary
  description: description
  parameters:
    - name: param1
      in: path
    - name: param2
      in: path
  get:
    summary: the get summary
    description: the get description
  post:
    summary: the post summary
    description: the post description
  servers:
    - url: the url
      description: the description
`

	item2 := `
/the/path:
  summary: summary2
  parameters:
    - name: param1
      in: path
    - name: param3
      in: query
  delete:
    summary: the delete summary
    description: the delete description
  post:
    summary: the post summary 2
    description: the post description 2
  servers:
    - url: the url 2
      description: the description 2
`

	var pathItems models.PathItems

	a := assert.New(t)
	a.NoError(yaml.Unmarshal([]byte(item1), &pathItems))
	a.NoError(yaml.Unmarshal([]byte(item2), &pathItems))

	a.Equal("ref", pathItems["/the/path"].Ref)
	a.Equal("summary2", pathItems["/the/path"].Summary)
	a.Equal("description", pathItems["/the/path"].Description)

	a.Equal("the get description", pathItems["/the/path"].Get.Description)
	a.Equal("the get summary", pathItems["/the/path"].Get.Summary)

	a.Equal("the post description 2", pathItems["/the/path"].Post.Description)
	a.Equal("the post summary 2", pathItems["/the/path"].Post.Summary)

	a.Equal("the delete description", pathItems["/the/path"].Delete.Description)
	a.Equal("the delete summary", pathItems["/the/path"].Delete.Summary)

	servers := pathItems["/the/path"].Servers
	a.Equal("the description", servers[0].Description)
	a.Equal("the url", servers[0].URL)
	a.Equal("the description 2", servers[1].Description)
	a.Equal("the url 2", servers[1].URL)

	parameters := pathItems["/the/path"].Parameters
	a.Equal("param1", parameters[0].Name)
	a.Equal("path", parameters[0].In)
	a.Equal("param2", parameters[1].Name)
	a.Equal("path", parameters[1].In)
	a.Equal("param3", parameters[2].Name)
	a.Equal("query", parameters[2].In)
}

func root() *models.Root {
	return &models.Root{
		OpenAPI: "3.0.2",
		Info: &models.Info{
			Title:          "something",
			Description:    "something",
			TermsOfService: "something",
			Version:        "something",
		},
		Servers: []*models.Server{
			{
				URL:         "something",
				Description: "something",
				Variables:   map[string]*models.ServerVariable{},
			},
		},
		Paths: map[string]*models.PathItem{
			"something": {
				Ref:         "something",
				Summary:     "something",
				Description: "something",
				Servers:     []*models.Server{},
				Parameters:  []*models.Parameter{},
			},
		},
		Components: &models.Components{
			Schemas:         map[string]*models.Schema{},
			Responses:       map[string]*models.Response{},
			Parameters:      map[string]*models.Parameter{},
			Examples:        map[string]*models.Example{},
			RequestBodies:   map[string]*models.RequestBody{},
			Headers:         map[string]*models.Header{},
			SecuritySchemes: map[string]*models.SecurityScheme{},
			Links:           map[string]*models.Link{},
			Callbacks:       map[string]*models.Callback{},
		},
		Security: []*models.SecurityRequirement{
			{Name: []string{"something"}},
		},
		Tags: []*models.Tag{
			{
				Name:                  "something",
				Description:           "something",
				ExternalDocumentation: nil,
			},
		},
		ExternalDocumentation: &models.ExternalDocumentation{
			Description: "something",
			URL:         "something",
		},
	}
}

func validateRoot(a *assert.Assertions, r *models.Root) {
	a.Equal("3.0.2", r.OpenAPI)
	a.Equal("something", r.Info.Title)
	a.Equal("something", r.Servers[0].URL)
	a.Equal("something", r.Paths["something"].Description)
	a.NotNil(r.Components)
	a.Equal("something", r.Security[0].Name[0])
	a.Equal("something", r.Tags[0].Name)
	a.Equal("something", r.ExternalDocumentation.URL)
}

func info() *models.Info {
	return &models.Info{
		Title:          "title",
		Description:    "description",
		TermsOfService: "tos",
		Contact: &models.Contact{
			Name:  "name",
			URL:   "url",
			Email: "email",
		},
		License: &models.License{
			Name: "name",
			URL:  "url",
		},
		Version: "version",
	}
}

func validateInfo(a *assert.Assertions, i *models.Info) {
	a.Equal("title", i.Title)
	a.Equal("description", i.Description)
	a.Equal("tos", i.TermsOfService)
	a.Equal("name", i.Contact.Name)
	a.Equal("url", i.Contact.URL)
	a.Equal("email", i.Contact.Email)
	a.Equal("name", i.License.Name)
	a.Equal("url", i.License.URL)
	a.Equal("version", i.Version)
}

func server() *models.Server {
	return &models.Server{
		URL:         "url",
		Description: "description",
		Variables: map[string]*models.ServerVariable{
			"key": {
				Enum:        []string{"enum"},
				Default:     "default",
				Description: "description",
			},
		},
	}
}

func validateServer(a *assert.Assertions, s *models.Server) {
	a.Equal("url", s.URL)
	a.Equal("description", s.Description)

	variable := s.Variables["key"]
	a.Equal("description", variable.Description)
	a.Equal("default", variable.Default)
	a.Equal("enum", variable.Enum[0])
}

func components() *models.Components {
	return &models.Components{
		Schemas: map[string]*models.Schema{
			"schema": {
				Nullable:   true,
				ReadOnly:   true,
				WriteOnly:  true,
				Deprecated: true,
			},
		},
		Responses: map[string]*models.Response{
			"response": {
				Description: "description",
			},
		},
		Parameters: map[string]*models.Parameter{
			"parameter": {
				Name:            "name",
				In:              "in",
				Description:     "description",
				Required:        true,
				Deprecated:      true,
				AllowEmptyValue: true,
				Style:           "style",
				Explode:         true,
				AllowReserved:   true,
				Example:         "example",
			},
		},
		Examples: map[string]*models.Example{
			"example": {
				Summary:       "summary",
				Description:   "description",
				Value:         "value",
				ExternalValue: "externalValue",
			},
		},
		RequestBodies: map[string]*models.RequestBody{
			"requestBody": {
				Description: "description",
				Required:    true,
			},
		},
		Headers: map[string]*models.Header{
			"header": {
				Description:     "description",
				Required:        true,
				Deprecated:      true,
				AllowEmptyValue: true,
				Style:           "style",
				Explode:         true,
				AllowReserved:   true,
				Example:         "example",
			},
		},
		SecuritySchemes: map[string]*models.SecurityScheme{
			"securitySchema": {
				Type:             "type",
				Description:      "description",
				Name:             "name",
				In:               "in",
				Scheme:           "scheme",
				BearerFormat:     "bearer",
				OpenIdConnectUrl: "url",
			},
		},
		Links: map[string]*models.Link{
			"link": {
				OperationRef: "ref",
				OperationID:  "id",
				Description:  "description",
			},
		},
		Callbacks: map[string]*models.Callback{
			"callback": {
				"pathItem": {
					Ref:         "ref",
					Summary:     "summary",
					Description: "description",
				},
			},
		},
	}
}

func validateComponents(a *assert.Assertions, c *models.Components) {
	a.Equal("description", c.Links["link"].Description)
	a.Equal("ref", c.Links["link"].OperationRef)
	a.Equal("id", c.Links["link"].OperationID)

	a.True(c.Schemas["schema"].Deprecated)
	a.True(c.Schemas["schema"].Nullable)
	a.True(c.Schemas["schema"].ReadOnly)
	a.True(c.Schemas["schema"].WriteOnly)

	a.Equal("ref", (*c.Callbacks["callback"])["pathItem"].Ref)
	a.Equal("description", (*c.Callbacks["callback"])["pathItem"].Description)

	a.Equal("description", c.Examples["example"].Description)
	a.Equal("summary", c.Examples["example"].Summary)
	a.Equal("externalValue", c.Examples["example"].ExternalValue)
	a.Equal("value", c.Examples["example"].Value)

	a.Equal("description", c.Headers["header"].Description)
	a.Equal("example", c.Headers["header"].Example)
	a.Equal("style", c.Headers["header"].Style)
	a.True(c.Headers["header"].Deprecated)
	a.True(c.Headers["header"].AllowEmptyValue)
	a.True(c.Headers["header"].AllowReserved)
	a.True(c.Headers["header"].Explode)
	a.True(c.Headers["header"].Required)

	a.Equal("description", c.Parameters["parameter"].Description)
	a.Equal("example", c.Parameters["parameter"].Example)
	a.Equal("style", c.Parameters["parameter"].Style)
	a.True(c.Parameters["parameter"].Deprecated)
	a.True(c.Parameters["parameter"].AllowEmptyValue)
	a.True(c.Parameters["parameter"].AllowReserved)
	a.True(c.Parameters["parameter"].Explode)
	a.True(c.Parameters["parameter"].Required)

	a.True(c.RequestBodies["requestBody"].Required)
	a.Equal("description", c.RequestBodies["requestBody"].Description)

	a.Equal("description", c.Responses["response"].Description)

	a.Equal("description", c.SecuritySchemes["securitySchema"].Description)
	a.Equal("name", c.SecuritySchemes["securitySchema"].Name)
	a.Equal("in", c.SecuritySchemes["securitySchema"].In)
	a.Equal("bearer", c.SecuritySchemes["securitySchema"].BearerFormat)
	a.Equal("url", c.SecuritySchemes["securitySchema"].OpenIdConnectUrl)
	a.Equal("scheme", c.SecuritySchemes["securitySchema"].Scheme)
	a.Equal("type", c.SecuritySchemes["securitySchema"].Type)
}
