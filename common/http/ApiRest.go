package apirest

import (
	"encoding/base64"

	// "github.com/sendgrid/rest"
	"go-firebase-gateway/common/http/rest"

	log "github.com/sirupsen/logrus"
)

type ApiRest struct {
}

var Client ApiRest

func (apiHttp *ApiRest) Get(baseUrl string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Get
	request.BaseURL = baseUrl
	log.Info("Get:  BaseURL: ", request.BaseURL, "; Params \n", request.QueryParams)
	return apiHttp.sendRequest(request)
}

func (apiHttp *ApiRest) Post(baseUrl string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Post
	request.BaseURL = baseUrl
	log.Info("Post: BaseURL: ", request.BaseURL, " request Body: ", string(request.Body))
	return apiHttp.sendRequest(request)
}
func (apiHttp *ApiRest) PostAuth(baseUrl, baseKey string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Post
	request.BaseURL = baseUrl
	headers := make(map[string]string)
	headers["Authorization"] = baseKey
	request.Headers = headers
	log.Info("Postv2: BaseURL: ", request.BaseURL, request.Headers, " request Body: ", string(request.Body))
	return apiHttp.sendRequest(request)
}

func (apiHttp *ApiRest) Patch(baseUrl string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Patch
	request.BaseURL = baseUrl
	log.Info("Patch: BaseURL: ", request.BaseURL, "; Params \n", string(request.Body))
	return apiHttp.sendRequest(request)
}

func (apiHttp *ApiRest) Put(baseUrl string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Put
	request.BaseURL = baseUrl
	log.Info("Put: BaseURL: ", request.BaseURL, "; Params \n", string(request.Body))
	return apiHttp.sendRequest(request)
}

func (apiHttp *ApiRest) Delete(baseUrl string, params map[string]interface{}) (rest.Response, error) {
	request := parseParams(params)
	request.Method = rest.Delete
	request.BaseURL = baseUrl
	log.Info("Delete: %v \n", request)
	return apiHttp.sendRequest(request)
}

func (apiHttp *ApiRest) sendRequest(request rest.Request) (rest.Response, error) {
	response, err := rest.Send(request)
	if err != nil {
		log.Error("ApiRest: %v", err)
		errResponse := rest.Response{}
		return errResponse, err
	}
	return *response, nil
}

func getDefaultHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return headers
}

func parseParams(params map[string]interface{}) rest.Request {
	request := rest.Request{}
	headersRaw := params["headers"]
	if headersRaw == nil {
		request.Headers = getDefaultHeaders()
	} else {
		headers, ok := headersRaw.(map[string]string)
		if !ok {
			log.Error("Headers is not passed as map[string]string")
			request.Headers = getDefaultHeaders()
		} else {
			request.Headers = headers
		}
	}
	queryParamsRaw := params["queryParams"]
	if queryParamsRaw == nil {
		request.QueryParams = make(map[string]string)
	} else {
		queryParams, ok := queryParamsRaw.(map[string]string)
		if !ok {
			log.Error("queryParams is not passed as map[string]string")
			request.QueryParams = make(map[string]string)
		} else {
			request.QueryParams = queryParams
		}

	}
	jsonBodyRaw := params["body"]
	if jsonBodyRaw == nil {
		request.Body = make([]byte, 0)
	} else {
		jsonBody, ok := jsonBodyRaw.(string)
		if !ok {
			log.Error("bodyParams is not passed as string")
			request.Body = make([]byte, 0)
		} else {
			var Body = []byte(jsonBody)
			request.Body = Body
		}
	}
	return request
}

func (apiHttp *ApiRest) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
