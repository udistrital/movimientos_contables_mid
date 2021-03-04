package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
)

func sendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Error(err)
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func getJsonTest(url string, target interface{}) (status int, err error) {
	r, err := http.Get(url)
	if err != nil {
		return r.StatusCode, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err.Error() + "error: " + url)
		}
	}()

	return r.StatusCode, json.NewDecoder(r.Body).Decode(target)
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}
