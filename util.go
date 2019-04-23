package harbor

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Query(content interface{}) (url.Values, error) {
	// v := reflect.ValueOf(content)
	// println(v.Kind())
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		return queryString(v.String())
	case reflect.Struct:
		return queryStruct(v.Interface())
	case reflect.Map:
		return queryMap(v.Interface())
	default:
	}
	return url.Values{}, nil
}

func queryStruct(content interface{}) (url.Values, error) {
	queryData := url.Values{}
	if marshalContent, err := json.Marshal(content); err != nil {
		return queryData, err
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			return queryData, err
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				var queryVal string
				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}
					queryVal = string(j)
				}
				queryData.Add(k, queryVal)
			}
		}
	}
	return queryData, nil
}

func queryString(content string) (url.Values, error) {
	queryData := url.Values{}
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			queryData.Add(k, v)
		}
	} else {
		if queryData, err := url.ParseQuery(content); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					queryData.Add(k, string(queryValue))
				}
			}
		} else {
			return queryData, err
		}
		// TODO: need to check correct format of 'field=val&field=val&...'
	}
	return queryData, nil
}

func queryMap(content interface{}) (url.Values, error) {
	return queryStruct(content)
}
