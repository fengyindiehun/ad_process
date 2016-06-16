package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"net/http"
	"time"
)

type httpProxy struct {
	urlGetFeedText     string
	urlGetCustType     string
	urlGetFeedTag      string
	urlGetCreativeType string
	timeout            int
}

func (proxy *httpProxy) requestUrl(url string) (string, error) {
	var respBody string
	client := http.Client{
		Timeout: time.Duration(proxy.timeout) * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return respBody, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respBody, err
	}
	return string(body), nil
}

type innerFeedMessage struct {
	Text string
}
type feedMessage struct {
	Statuses []innerFeedMessage
}

func decodeFeed(response string) (text string, err error) {
	var m feedMessage
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return
	}
	if len(m.Statuses) < 1 {
		err = fmt.Errorf("text parse error, response:%s", response)
		return
	}
	text = m.Statuses[0].Text
	return
}

func (proxy *httpProxy) getFeedText(feedId string) (text string, err error) {
	url := proxy.urlGetFeedText + feedId
	response, err := proxy.requestUrl(url)
	if err != nil {
		return
	}
	text, err = decodeFeed(response)
	return
}

type custMessage struct {
	Verified      bool
	Verified_type int
}

func decodeCustInfo(response string) (custType string, err error) {
	var m custMessage
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return
	}
	if verified, verifiedType := m.Verified, m.Verified_type; verified == false {
		custType = "3"
	} else {
		if verifiedType == 0 {
			custType = "1"
		} else if verifiedType > 0 {
			custType = "2"
		} else {
			err = fmt.Errorf("custType parse error, response:%s", response)
		}
	}
	return
}

//custType: 1<=>Orange 2<=>Blue 3<=>Ordinary
//doc: http://wiki.intra.sina.com.cn/pages/viewpage.action?pageId=102040429
func (proxy *httpProxy) getCustType(custId string) (custType string, err error) {
	url := proxy.urlGetCustType + custId
	response, err := proxy.requestUrl(url)
	if err != nil {
		return
	}
	custType, err = decodeCustInfo(response)
	return
}

func (proxy *httpProxy) getFeedTag(feeId, custId string, feedText string) (tag string, err error) {
	//url := proxy.urlGetFeedTag + custId
	//tag, err = proxy.requestUrl(url)
	return
}

func (proxy *httpProxy) getCreativeType(feedId string) (creativeType string, err error) {
	//url := proxy.urlGetCreativeType + feedId
	//creativeType, err = proxy.requestUrl(url)
	return
}

func newHttpProxy(section *ini.Section) (*httpProxy, error) {
	urlGetFeedText := section.Key(HTTP_SECTION_GET_FEEDTEXT).String()
	if urlGetFeedText == "" {
		return nil, fmt.Errorf("urlGetFeedText parse error")
	}

	urlGetCustType := section.Key(HTTP_SECTION_GET_CUSTTYPE).String()
	if urlGetCustType == "" {
		return nil, fmt.Errorf("urlGetCustType parse error")
	}

	urlGetFeedTag := section.Key(HTTP_SECTION_GET_FEEDTAG).String()
	if urlGetFeedTag == "" {
		return nil, fmt.Errorf("urlGetFeedTag parse error")
	}

	urlGetCreativeType := section.Key(HTTP_SECTION_GET_CREATIVETYPE).String()
	if urlGetCreativeType == "" {
		return nil, fmt.Errorf("urlGetCreativeType parse error")
	}

	timeout, err := section.Key(HTTP_SECTION_TIMEOUT).Int()
	if err != nil {
		return nil, fmt.Errorf("timeout parse error")
	}

	return &httpProxy{
		urlGetFeedText:     urlGetFeedText,
		urlGetCustType:     urlGetCustType,
		urlGetFeedTag:      urlGetFeedTag,
		urlGetCreativeType: urlGetCreativeType,
		timeout:            timeout,
	}, nil
}
