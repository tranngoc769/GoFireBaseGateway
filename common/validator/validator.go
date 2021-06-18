package validator

import (
	"fmt"
	"go-firebase-gateway/common/response"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
)

func getSchema(schema string) ([]byte, error) {
	return ioutil.ReadFile("./common/schema/" + schema)
}

func GetSchema(schema string) ([]byte, error) {
	return getSchema(schema)
}

func CheckSchema(schema string, value interface{}) (int, interface{}) {
	schemaChecker, err := getSchema(schema)
	if err != nil {
		return response.NotFoundMsg(schema + " is not existed")
	}
	schemaLoader := gojsonschema.NewStringLoader(fmt.Sprintf("%v", string(schemaChecker)))
	documentLoader := gojsonschema.NewGoLoader(value)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if result.Valid() {
		return response.NewOKResponse("")
	} else {
		var errArr []map[string]interface{}
		for _, msg := range result.Errors() {
			errMsg := map[string]interface{}{
				msg.Field(): msg.Description(),
			}
			errArr = append(errArr, errMsg)
		}
		return response.NewErrorResponse(400, errArr)
	}
}
