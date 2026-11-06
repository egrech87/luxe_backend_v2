package apiUtils

import (
	"Quote/db/daoPostgres/daoQuery"
	"Quote/utils/arrayUtils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Notes when reading multiple values, that is, an array:
// - It is important that the values of the array are separated by the VALUES_SEPARATOR constant
// - The defaultType parameter must be one of the HTTP_DEFAULT_?_ARRAY constants, such as HTTP_DEFAULT_NIL_ARRAY and HTTP_DEFAULT_EMPTY_ARRAY

const (
	VALUES_SEPARATOR_QUERY  string = "|"
	VALUES_SEPARATOR_HEADER string = ","
	VALUES_SEPARATOR_PATH   string = ","
	VALUES_SEPARATOR_BODY   string = ","

	HTTP_QUERY int8 = iota
	HTTP_PATH
	HTTP_HEADER
	HTTP_BODY

	HTTP_DEFAULT_NIL_ARRAY
	HTTP_DEFAULT_EMPTY_ARRAY
)

//#region READ QUERY PARAMETERS

func HttpQueryString(c *gin.Context, param string, isRequired bool) (*string, error) {
	if !c.Request.URL.Query().Has(param) {
		if isRequired {
			return nil, fmt.Errorf("'%s' query parameter is missing", param)
		} else {
			return nil, nil
		}
	}

	value := c.Request.URL.Query().Get(param)
	return &value, nil
}

func HttpQueryInt(c *gin.Context, param string, isRequired bool) (*int, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	value, err := strconv.Atoi(*valueS)
	if err != nil {
		return nil, fmt.Errorf("'%s' query parameter must be a whole number: %s", param, err.Error())
	}
	return &value, nil
}

func HttpQueryUInt64(c *gin.Context, param string, isRequired bool) (*uint64, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	value, err := strconv.ParseUint(*valueS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("'%s' query parameter must be a positive whole number: %s", param, err.Error())
	}
	return &value, nil
}

// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Any other value returns an error.
func HttpQueryBool(c *gin.Context, param string, isRequired bool) (*bool, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	value, err := strconv.ParseBool(*valueS)
	if err != nil {
		return nil, fmt.Errorf("'%s' query parameter must be either true or false: %s", param, err.Error())
	}
	return &value, nil
}

// Reads a date in format yyyy-MM-dd and converts it to Unix Epoch timestamp
func HttpQueryFromDate(c *gin.Context, param string, isRequired bool) (*time.Time, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	fromTime, err := time.Parse(time.RFC3339, *valueS+"T00:00:00+00:00")
	if err != nil {
		return nil, fmt.Errorf("incorrect '%s' date format. must be in the format yyyy-MM-dd", param)
	}

	return &fromTime, nil
}

// Reads a date in format yyyy-mm-dd and converts it to Unix Epoch timestamp
func HttpQueryToDate(c *gin.Context, param string, isRequired bool) (*time.Time, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	toTime, err := time.Parse(time.RFC3339, *valueS+"T23:59:59.999999999+00:00")
	if err != nil {
		return nil, fmt.Errorf("incorrect '%s' date format. must be in the format yyyy-MM-dd", param)
	}

	return &toTime, nil
}

func HttpQueryInt_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]int, error) {
	return httpInt_Array(c, param, isRequired, defaultType, HTTP_QUERY)
}

func HttpQueryUInt64_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]uint64, error) {
	return httpUInt64_Array(c, param, isRequired, defaultType, HTTP_QUERY)
}

func HttpQueryString_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]string, error) {
	return httpString_Array(c, param, isRequired, defaultType, HTTP_QUERY)
}

// Reads a date-time in format yyyy-MM-ddTHH:mm:ssZ and converts it to Unix Epoch timestamp
func HttpQueryFromDateTime(c *gin.Context, param string, isRequired bool) (*time.Time, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	fromTime, err := time.Parse(time.RFC3339, *valueS)
	if err != nil {
		return nil, fmt.Errorf("incorrect '%s' date-time format. must be in the format yyyy-MM-ddTHH:mm:ssZ", param)
	}

	return &fromTime, nil
}

// Reads a date-time in format yyyy-MM-ddTHH:mm:ssZ and converts it to Unix Epoch timestamp
func HttpQueryToDateTime(c *gin.Context, param string, isRequired bool) (*time.Time, error) {
	valueS, err := HttpQueryString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	toTime, err := time.Parse(time.RFC3339, *valueS)
	if err != nil {
		return nil, fmt.Errorf("incorrect '%s' date-time format. must be in the format yyyy-MM-ddTHH:mm:ssZ", param)
	}

	// To include everything up to the given second. Ignore this if nanoseconds are already provided in the parameter as part of the 'to' value
	if !strings.Contains(*valueS, ".") {
		toTime = toTime.Add(time.Nanosecond * 999999999)
	}

	return &toTime, nil
}

func HttpQueryPagination(c *gin.Context) (*daoQuery.Pagination, error) {
	p := daoQuery.NewPagination()

	sort := ""
	sortBy := ""
	sortDirection := "desc"

	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "pagesize":
			limit, err := strconv.Atoi(queryValue)
			if err != nil {
				return p, errors.Errorf("'pagesize' query parameter is not whole number: %v", err)
			}
			p.SetLimit(limit)
		case "page":
			page, err := strconv.Atoi(queryValue)
			if err != nil {
				return p, errors.Errorf("'page' query parameter is not whole number: %v", err)
			}
			p.SetPage(page)
		case "sort":
			sort = queryValue
		case "sort_by":
			sortBy = queryValue
		case "sort_direction":
			sortDirection = queryValue
		}
	}

	if sort == "" {
		if sortBy != "" {
			sort = sortBy + " " + sortDirection
		}
	}

	if sort != "" {
		p.SetSort(sort)
	}

	return p, nil
}

//#endregion

//#region READ PATH PARAMETERS

func HttpPathString(c *gin.Context, param string) (*string, error) {
	value := c.Param(param)
	return &value, nil
}

func HttpPathInt(c *gin.Context, param string) (*int, error) {
	valueS := c.Param(param)
	value, err := strconv.Atoi(valueS)
	if err != nil {
		return nil, fmt.Errorf("'"+param+"' path parameter must be a whole number: %v", err.Error())
	}
	return &value, nil
}

func HttpPathUInt64(c *gin.Context, param string) (*uint64, error) {
	valueS := c.Param(param)
	value, err := strconv.ParseUint(valueS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("'"+param+"' path parameter must be a positive whole number: %v", err.Error())
	}
	return &value, nil
}

func HttpPathFloat64(c *gin.Context, param string) (*float64, error) {
	valueS := c.Param(param)
	value, err := strconv.ParseFloat(valueS, 64)
	if err != nil {
		return nil, fmt.Errorf("'"+param+"' path parameter must be a number: %v", err.Error())
	}
	return &value, nil
}

func HttpPathInt_Array(c *gin.Context, param string, defaultType int8) (*[]int, error) {
	return httpInt_Array(c, param, true, defaultType, HTTP_PATH)
}

func HttpPathUInt64_Array(c *gin.Context, param string, defaultType int8) (*[]uint64, error) {
	return httpUInt64_Array(c, param, true, defaultType, HTTP_PATH)
}

func HttpPathString_Array(c *gin.Context, param string, defaultType int8) (*[]string, error) {
	return httpString_Array(c, param, true, defaultType, HTTP_PATH)
}

//#endregion

// #region READ HEADER PARAMETERS
func HttpHeaderString(c *gin.Context, param string, isRequired bool) (*string, error) {
	value, exists := c.Request.Header[param]
	if !exists {
		if isRequired {
			return nil, fmt.Errorf("'" + param + "' header parameter is missing")
		} else {
			return nil, nil
		}
	}
	return &value[0], nil
}

func HttpHeaderInt_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]int, error) {
	return httpInt_Array(c, param, isRequired, defaultType, HTTP_HEADER)
}

func HttpHeaderUInt64_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]uint64, error) {
	return httpUInt64_Array(c, param, isRequired, defaultType, HTTP_HEADER)
}

func HttpHeaderString_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]string, error) {
	return httpString_Array(c, param, isRequired, defaultType, HTTP_HEADER)
}

func HttpHeaderJson(c *gin.Context, param string, isRequired bool) (map[string]any, error) {
	valueS, err := HttpHeaderString(c, param, isRequired)
	if valueS == nil {
		return nil, err
	}

	value := map[string]any{}
	err = json.Unmarshal([]byte(*valueS), &value)
	return value, err
}

func HttpHeaderSelector(c *gin.Context) (selector *daoQuery.Selector, err error) {
	jsonMap, err := HttpHeaderJson(c, "Selector", false)
	if err != nil {
		return selector, errors.Errorf("failed to decode header Selector: %s", err.Error())
	}

	selector = &daoQuery.Selector{}
	err = selector.FillFromJson(jsonMap)
	if err != nil {
		return selector, errors.Errorf("failed to decode header Selector: %s", err.Error())
	}

	return selector, err
}

//#endregion

// #region READ BODY PARAMETERS

func HttpBodyJson(c *gin.Context, isRequired bool) (map[string]any, error) {
	value := map[string]any{}

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return value, errors.Errorf("failed to decode request body: %s", err.Error())
	}
	if len(requestBody) == 0 && !isRequired {
		return value, nil
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	err = json.NewDecoder(requestBodyBuffer).Decode(&value)
	if err != nil {
		return value, errors.Errorf("failed to decode request body: %s", err.Error())
	}

	return value, nil
}
func HttpBodyStruct(c *gin.Context, structure any) error {
	fmt.Println("Type of variable1:", reflect.TypeOf(structure))
	return json.NewDecoder(c.Request.Body).Decode(&structure)
}

func HttpBodyInt_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]int, error) {
	return httpInt_Array(c, param, isRequired, defaultType, HTTP_BODY)
}
func HttpBodyUInt64_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]uint64, error) {
	return httpUInt64_Array(c, param, isRequired, defaultType, HTTP_BODY)
}
func HttpBodyString_Array(c *gin.Context, param string, isRequired bool, defaultType int8) (*[]string, error) {
	return httpString_Array(c, param, isRequired, defaultType, HTTP_BODY)
}

//#endregion

// #region READ MultipartForm

func HttpMultipartForm(c *gin.Context, param string) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []*multipart.FileHeader{}, errors.Errorf(fmt.Sprintf("failed to get the file(s): %v", err.Error()))
	}

	filesInput := form.File[param]

	return filesInput, nil
}

// #endregion

//#region private Utility functions

func httpInt_Array(c *gin.Context, param string, isRequired bool, defaultType int8, httpParameterType int8) (*[]int, error) {
	defaultValue := &[]int{}
	if defaultType == HTTP_DEFAULT_NIL_ARRAY {
		defaultValue = nil
	}

	allValuesStr, err := httpString_Array(c, param, isRequired, defaultType, httpParameterType)
	if err != nil || allValuesStr == nil {
		return defaultValue, err
	}

	values, err := arrayUtils.StringArrayToIntArray(*allValuesStr, arrayUtils.ALL_ITEMS)
	return &values, err
}

func httpUInt64_Array(c *gin.Context, param string, isRequired bool, defaultType int8, httpParameterType int8) (*[]uint64, error) {
	defaultValue := &[]uint64{}
	if defaultType == HTTP_DEFAULT_NIL_ARRAY {
		defaultValue = nil
	}

	allValuesStr, err := httpString_Array(c, param, isRequired, defaultType, httpParameterType)
	if err != nil || allValuesStr == nil {
		return defaultValue, err
	}

	values, err := arrayUtils.StringArrayToUint64Array(*allValuesStr, arrayUtils.ALL_ITEMS)
	return &values, err
}

func httpString_Array(c *gin.Context, param string, isRequired bool, defaultType int8, httpParameterType int8) (*[]string, error) {
	defaultValue := &[]string{}
	if defaultType == HTTP_DEFAULT_NIL_ARRAY {
		defaultValue = nil
	}

	var allValues string
	var httpTypeName string
	var separator string
	switch httpParameterType {
	case HTTP_QUERY:
		allValues = c.Request.URL.Query().Get(param)
		httpTypeName = "query"
		separator = VALUES_SEPARATOR_QUERY
	case HTTP_PATH:
		allValues = c.Param(param)
		httpTypeName = "path"
		separator = VALUES_SEPARATOR_PATH
	case HTTP_HEADER:
		allValues = c.GetHeader(param)
		httpTypeName = "header"
		separator = VALUES_SEPARATOR_HEADER
	case HTTP_BODY:
		type ArrayBody struct {
			ArrayInput string
		}
		body := ArrayBody{}
		bodyStruct := reflect.TypeOf(body)
		structFields := []reflect.StructField{bodyStruct.Field(0)}
		structFields[0].Tag = reflect.StructTag(`json:"` + param + `"`)
		newBodyStruct := reflect.StructOf(structFields)
		newBodyValue := reflect.New(newBodyStruct).Elem()
		newBody := newBodyValue.Addr().Interface()
		json.NewDecoder(c.Request.Body).Decode(newBody)
		bodyStr := fmt.Sprintf("%v", newBody)
		bodyStr = strings.Trim(bodyStr, "&{")
		bodyStr = strings.Trim(bodyStr, "}")
		allValues = bodyStr
		httpTypeName = "body"
		separator = VALUES_SEPARATOR_BODY
	default:
		return defaultValue, fmt.Errorf("failed to read '%v' parameter", param)
	}

	if allValues == "" {
		if isRequired {
			return defaultValue, fmt.Errorf("'%v' %v parameter is missing", param, httpTypeName)
		} else {
			return defaultValue, nil
		}
	}

	var err error
	allValues, err = url.QueryUnescape(allValues)
	if err != nil {
		return defaultValue, fmt.Errorf("'%v' %v parameter has invalid characters: %v", param, httpTypeName, err.Error())
	}
	allValues = strings.Trim(allValues, " ")

	valueS_array := strings.Split(allValues, separator)

	final_values := []string{}
	for _, valueS := range valueS_array {
		trimmedValueS := strings.Trim(valueS, " ")
		if len(trimmedValueS) > 0 {
			final_values = append(final_values, trimmedValueS)
		}
	}

	return &final_values, nil
}

//#endregion
