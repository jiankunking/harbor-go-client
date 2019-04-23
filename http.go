package harbor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jiankunking/errwrap"
)

func httpGet(c *Client, queryData url.Values, subPath string) ([]byte, *http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, subPath)
	//fmt.Println(url)
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, nil, errReq
	}

	if queryData != nil {
		// Add all querystring from Query func
		q := req.URL.Query()
		for k, v := range queryData {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	// Add basic auth
	if c.BasicAuth != struct{ Username, Password string }{} {
		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if errReq != nil {
		return nil, nil, errReq
	}
	resp, errDo := c.Client.Do(req)
	if errDo != nil {
		return nil, resp, errDo
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp, errwrap.WrapString("request status code exception : " + strconv.Itoa(resp.StatusCode))
	}

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, resp, errRead
	}
	return body, resp, nil
}

func httpDelete(c *Client, subPath string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, subPath)
	req, errReq := http.NewRequest(http.MethodDelete, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	// Add basic auth
	if c.BasicAuth != struct{ Username, Password string }{} {
		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}
	resp, errDo := c.Client.Do(req)
	if errDo != nil {
		return resp, errDo
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp, errwrap.WrapString("request status code exception : " + strconv.Itoa(resp.StatusCode))
	}
	return resp, nil
}

//func httpPost(c *Client, subPath string, body []byte) (*http.Response, error) {
//	url := fmt.Sprintf("%s/%s", c.BaseURL, subPath)
//	req, errReq := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
//	if errReq != nil {
//		return nil, errReq
//	}
//	// Add basic auth
//	if c.BasicAuth != struct{ Username, Password string }{} {
//		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
//	}
//	resp, errDo := c.Client.Do(req)
//	if errDo != nil {
//		return resp, errDo
//	}
//	defer resp.Body.Close()
//	return resp, nil
//}
//func httpPut(c *Client, subPath string, body []byte) (*http.Response, error) {
//	url := fmt.Sprintf("%s/%s", c.BaseURL, subPath)
//	req, errReq := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//	if errReq != nil {
//		return nil, errReq
//	}
//	// Add basic auth
//	if c.BasicAuth != struct{ Username, Password string }{} {
//		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
//	}
//	resp, errDo := c.Client.Do(req)
//	if errDo != nil {
//		return resp, errDo
//	}
//	defer resp.Body.Close()
//	return resp, nil
//}
