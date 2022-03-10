package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func CreateAccount(input *AccountData) (*AccountData, error) {
	req := Request{
		Data: input,
	}

	body, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(buildUrl("/organisation/accounts"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected response code: %v\nbody: %v", resp.StatusCode, string(bytes))
	}

	var ret Response
	err = json.Unmarshal(bytes, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Data, nil
}

func buildUrl(path string) string {
	baseUrl, ok := os.LookupEnv("BASEURL")
	if !ok {
		panic("BASEURL variable must be set")
	}

	if strings.HasSuffix(baseUrl, "/") {
		if strings.HasPrefix(path, "/") {
			return fmt.Sprintf("%v%v", baseUrl[:len(baseUrl)-1], path)
		}
		return fmt.Sprintf("%v%v", baseUrl, path)
	}

	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("%v%v", baseUrl, path)
	}
	return fmt.Sprintf("%v/%v", baseUrl, path)
}

type Request struct {
	Data *AccountData `json:"data"`
}

type Response struct {
	Data *AccountData `json:"data"`
}
