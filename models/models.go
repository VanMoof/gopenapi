package models

import "gopkg.in/yaml.v3"

type Root struct {
	OpenAPI               string                 `json:"openapi" yaml:"openapi"`
	Info                  *Info                  `json:"info" yaml:"info"`
	Servers               []*Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths                 PathItems              `json:"paths" yaml:"paths"`
	Components            *Components            `json:"components,omitempty" yaml:"components,omitempty"`
	Security              []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Tags                  []*Tag                 `json:"tag,omitempty" yaml:"tag,omitempty"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type PathItems map[string]*PathItem

func (n *PathItems) UnmarshalYAML(value *yaml.Node) error {
	if *n == nil {
		*n = map[string]*PathItem{}
	}

	decodedItems := map[string]*PathItem{}
	if err := value.Decode(&decodedItems); err != nil {
		return err
	}

	for decodedPath, decodedPathItem := range decodedItems {
		if existingPathItem, ok := (*n)[decodedPath]; ok {
			existingPathItem.merge(decodedPathItem)
		} else {
			(*n)[decodedPath] = decodedPathItem
		}
	}
	return nil
}

type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

type PathItem struct {
	Ref         string       `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string       `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation   `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation   `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation   `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation   `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation   `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation   `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation   `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation   `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*Server    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (p *PathItem) merge(other *PathItem) {
	if other.Ref != "" {
		p.Ref = other.Ref
	}
	if other.Summary != "" {
		p.Summary = other.Summary
	}
	if other.Description != "" {
		p.Description = other.Description
	}
	if other.Get != nil {
		p.Get = other.Get
	}
	if other.Put != nil {
		p.Put = other.Put
	}
	if other.Post != nil {
		p.Post = other.Post
	}
	if other.Delete != nil {
		p.Delete = other.Delete
	}
	if other.Options != nil {
		p.Options = other.Options
	}
	if other.Head != nil {
		p.Head = other.Head
	}
	if other.Patch != nil {
		p.Patch = other.Patch
	}
	if other.Trace != nil {
		p.Trace = other.Trace
	}
	if len(other.Servers) > 0 {
		p.Servers = append(p.Servers, other.Servers...)
	}
	if len(other.Parameters) > 0 {
		existingParams := map[string]interface{}{}
		for _, param := range p.Parameters {
			existingParams[param.Name+param.In] = true
		}

		for _, otherParam := range other.Parameters {
			if _, ok := existingParams[otherParam.Name+otherParam.In]; !ok {
				p.Parameters = append(p.Parameters, otherParam)
			}
		}
	}
}

type Operation struct {
	Tags                  []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary               string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description           string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID           string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters            []*Parameter           `json:"parameter,omitempty" yaml:"parameter,omitempty"`
	RequestBody           *RequestBody           `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses             map[string]*Response   `json:"responses" yaml:"responses"`
	Callbacks             map[string]*Callback   `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated            bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security              []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	Servers               []*Server              `json:"servers,omitempty" yaml:"servers,omitempty"`
}

type Parameter struct {
	Name            string                `json:"name" yaml:"name"`
	In              string                `json:"in" yaml:"in"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required" yaml:"required"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         bool                  `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}           `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content" yaml:"content"`
}

type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content" yaml:"content"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
}

type MediaType struct {
	Schema   *Schema              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}          `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Example  `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

type Header struct {
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required" yaml:"required"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         bool                  `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}           `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content" yaml:"content"`
}

type Encoding struct {
	ContentType   string             `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string             `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool               `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool               `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

type Response struct {
	Ref         string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Headers     map[string]*Header    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty" yaml:"links,omitempty"`
}

type Link struct {
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string                 `json:"description" yaml:"description"`
	Server       *Server                `json:"server,omitempty" yaml:"server,omitempty"`
}

type Example struct {
	Summary       string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty" yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

type Tag struct {
	Name                  string                 `json:"name" yaml:"name"`
	Description           string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type SecurityRequirement struct {
	Name []string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ExternalDocumentation struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

type Callback map[string]*PathItem

type Schema struct {
	Ref                   string                 `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Nullable              bool                   `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Discriminator         *Discriminator         `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	ReadOnly              bool                   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly             bool                   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	XML                   *XML                   `json:"xml,omitempty" yaml:"xml,omitempty"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Example               interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated            bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	Type                 string             `json:"type,omitempty" yaml:"type,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  []*Schema          `json:"not,omitempty" yaml:"not,omitempty"`
	Items                *Schema            `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string             `json:"description,omitempty" yaml:"description,omitempty"`
	Default              interface{}        `json:"default,omitempty" yaml:"default,omitempty"`
	Format               string             `json:"format,omitempty" yaml:"format,omitempty"`
}

type XML struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

type SecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OathFlowObject `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OathFlowObject `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OathFlowObject `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OathFlowObject `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

type OathFlowObject struct {
	AuthorizationURL string            `json:"authorizationUrl" yaml:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl" yaml:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}
