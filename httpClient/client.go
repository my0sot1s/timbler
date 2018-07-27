package httpClient

import (
	"io/ioutil"
	"net/http"
)

func Get(endpoint string) ([]byte, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

func Post(endpoint string, forms map[string][]string) ([]byte, error) {
	res, err := http.PostForm(endpoint, forms)
	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}
