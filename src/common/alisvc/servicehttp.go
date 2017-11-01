package alisvc

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func SvcPost(urlstr string, body map[string]string, accessId string, accessKey string) (string, error) {
	method := "POST"
	accept := "application/json"
	content_type := "application/json"
	urlobj, _ := url.Parse(urlstr)
	path := urlobj.Path
	date := time.Now().UTC().Format("Mon Jan 2 15:04:05 -0700 MST 2006")

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(&body)
	bodystr := buf.String()
	h := md5.New()
	h.Write(buf.Bytes())
	cipherStr := h.Sum(nil)
	bodyMd5 := base64.StdEncoding.EncodeToString(cipherStr)

	stringToSign := method + "\n" + accept + "\n" + bodyMd5 + "\n" + content_type + "\n" + date + "\n" + path
	// HMAC-SHA1
	mac := hmac.New(sha1.New, []byte(accessKey))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// authorization header
	authHeader := "Dataplus " + accessId + ":" + signature

	req, _ := http.NewRequest("POST", urlstr, strings.NewReader(bodystr))

	req.Header.Set("accept", accept)
	req.Header.Set("content-type", content_type)
	req.Header.Set("date", date)
	req.Header.Set("Authorization", authHeader)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {

		return "", err
	}

	defer resp.Body.Close()
	res, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}
	return string(res), nil

}
func SvcPost_Old(urlstr string, body map[string]string, accessId string, accessKey string) (string, error) {

	if _, ok := body["Format"]; !ok {
		body["Format"] = "json"
	}
	if _, ok := body["SignatureMethod"]; !ok {
		body["SignatureMethod"] = "HMAC-SHA1"
	}
	if _, ok := body["SignatureVersion"]; !ok {
		body["SignatureVersion"] = "1.0"
	}
	if _, ok := body["SignatureNonce"]; !ok {
		body["SignatureNonce"] = "1.0"
	}

	localutc, _ := time.LoadLocation("UTC")
	body["Timestamp"] = time.Now().In(localutc).Format("2006-01-02T15:04:05Z")
	body["SignatureNonce"] = fmt.Sprintf("%d%d", rand.Intn(0xFFFFFFFF), time.Now().Unix())
	body["AccessKeyId"] = accessId
	delete(body, "Signature")
	body["Signature"] = AliYunSign_Old(body, accessKey)

	lines := make([]string, len(body))
	idx := 0
	for key, val := range body {
		if len(val) == 0 {
			lines[idx] = url.QueryEscape(key)
		} else {
			lines[idx] = url.QueryEscape(key) + "=" + url.QueryEscape(val)
		}
		idx++
	}
	/*
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.Encode(&body)
		//return buf.String()
		http.Post
		resp, err := http.Post(urlstr,
			"application/x-www-form-urlencoded",
			strings.NewReader(buf.String()))
		if err != nil {
			fmt.Println(err)
		}*/
	//return buf.String()

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}
	//resp, err := client.Get("https://example.com")

	resp, err := client.Post(urlstr,
		"application/x-www-form-urlencoded",
		strings.NewReader(strings.Join(lines, "&")))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	res, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}
	return string(res), nil
}
func AliYunSign_Old(body map[string]string, accessKey string) string {
	//testsign := "POST&%2F&AccessKeyId%3DlzDLASvYm0becL4c%26Action%3DImageDetection%26Async%3Dfalse%26Format%3Djson%26ImageUrl%3Ddemo/g/logo.png%26RegionId%3Dcn-hangzhou%26Scene%3Dporn%26SignatureMethod%3DHMAC-SHA1%26SignatureNonce%3D18042893831477980729%26SignatureVersion%3D1.0%26Timestamp%3D2016-11-01T06%253A12%253A09Z%26Version%3D2016-08-01"
	//         POST&%2F&AccessKeyId%3DD2MTsDxKCp5M5R90%26Action%3DImageDetection%26Async%3Dfalse%26Format%3Djson%26ImageUrl.1%3g%252Flogo.png%26                       Scene.1%3Dporn%26SignatureMethod%3DHMAC-SHA1%26SignatureNonce%3D14325185151507709355%26SignatureVersion%3D1.0%26Timestamp%3D2017-10-11T08%253A09%253A15Z%26Version%3D2016-08-01
	testkey := accessKey + "&"

	b := bytes.Buffer{}
	b.WriteString("POST&%2F&")

	sorted_keys := make([]string, 0)
	for k, _ := range body {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	lines := make([]string, len(sorted_keys))
	for idx, k := range sorted_keys {
		val := body[k]
		if len(val) == 0 {
			lines[idx] = k
		} else {
			lines[idx] = k + "=" + url.QueryEscape(val)
		}
	}
	b.WriteString(url.QueryEscape(strings.Join(lines, "&")))

	//hmac ,use sha1
	mac := hmac.New(sha1.New, []byte(testkey))
	mac.Write(b.Bytes())
	encodeString := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return encodeString
}
