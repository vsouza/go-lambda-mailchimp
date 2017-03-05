package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	apex "github.com/apex/go-apex"
)

type subscriber struct {
	Email       string                 `json:"email_address"`
	Status      string                 `json:"status"`
	MergeFields map[string]interface{} `json:"merge_fields"`
}

func subscribe(subs *subscriber) error {

	subscriber := &subscriber{
		Email:       subs.Email,
		Status:      "subscribed",
		MergeFields: subs.MergeFields,
	}

	json, err := json.Marshal(subscriber)
	if err != nil {
		return err
	}

	var url bytes.Buffer
	url.WriteString("https://us12.api.mailchimp.com/3.0/lists/")
	url.WriteString("911e64fe78")
	url.WriteString("/members")
	fmt.Println(url.String())
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(json))
	req.SetBasicAuth("anystring", "098119cc4c622d98f02ab829803e9ace-us12")
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var s subscriber
		if err := json.Unmarshal(event, &s); err != nil {
			return nil, err
		}
		if err := subscribe(&s); err != nil {
			return nil, err
		}
		return map[string]interface{}{"status": "success"}, nil
	})
}
