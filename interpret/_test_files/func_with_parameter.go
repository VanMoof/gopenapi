// +build testResource

package _test_files

/*
gopenapi:parameter
in: path
required: true
content:
  text/plain:
    example: some text
*/
const ConstParamName = "constParamName"

/*
gopenapi:parameter
in: query
required: true
content:
  text/plain:
    example: some text
*/
var VarParamName = "varParamName"
