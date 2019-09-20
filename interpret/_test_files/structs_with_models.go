// +build testResource

package _test_files

import "time"

//gopenapi:objectSchema
type RootModel struct {
	IntField    int64       `json:"intField"`
	StringField string      `json:"stringField"`
	SubModels   []*SubModel `json:"subModels"`
}

//gopenapi:objectSchema
type SubModel struct {
	FloatField  float64                 `json:"floatField"`
	SubSubModel map[string]*SubSubModel `json:"subSubModel"`
}

//gopenapi:objectSchema
type SubSubModel struct {
	BoolField bool        `json:"boolField"`
	Aliased   AliasedSubs `json:"aliased"`
}

type IgnoredModel struct {
}

//gopenapi:objectSchema
type AliasedSubs []*AliasedSub

//gopenapi:objectSchema
type AliasedSub struct {
	IgnoredField string `json:"-"`
	TimeField    time.Time
}
