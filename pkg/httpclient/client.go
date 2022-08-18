package httpclient

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeHttpReq(apiKey string, req *http.Request) []byte {
	req.Header.Add("Authorization", "Basic "+basicAuth(apiKey, "X"))
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes
}
