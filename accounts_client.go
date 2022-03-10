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

// Send an account creation request to the API, returning the resulting account object.
// Callers should use the returned object and not the one given to this function to ensure
// any fields populated by the API are present in the following code.
func CreateAccount(input *AccountData) (*AccountData, error) {
	body := marshal(input)

	resp, err := http.Post(buildUrl("/organisation/accounts"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, &ConnectionError{err}
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ConnectionError{err} // not really a connection error, but doesn't belong in client or server side categories either
	}

	if resp.StatusCode > 299 && resp.StatusCode < 500 {
		return nil, &ClientError{StatusCode: resp.StatusCode, Message: string(bytes)}
	}
	if resp.StatusCode >= 500 {
		return nil, &ServerError{StatusCode: resp.StatusCode, Message: string(bytes)}
	}

	ret := unmarshal(bytes)
	return &ret, nil
}

// Fetch the account with the given uuid. If the account doesn't exist the function will return (nil, nil).
func FetchAccount(accountId string) (*AccountData, error) {
	path := fmt.Sprintf("/organisation/accounts/%v", accountId)

	resp, err := http.Get(buildUrl(path))
	if err != nil {
		return nil, &ConnectionError{err}
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ConnectionError{err}
	}

	if resp.StatusCode == 200 {
		ret := unmarshal(bytes)
		return &ret, nil
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	if resp.StatusCode >= 500 {
		return nil, &ServerError{StatusCode: resp.StatusCode, Message: string(bytes)}
	}

	return nil, fmt.Errorf("unexpected status code %d: %v", resp.StatusCode, string(bytes))
}

// builds a complete url from the env var BASEURL + the given path.
// trailing / in baseUrl and leading / in path will be taken into account.
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

// marshal the given AccountData into a Request, as json, in a []byte.
// will panic if marshalling fails.
func marshal(input *AccountData) []byte {
	req := Request{
		Data: input,
	}

	body, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}

	return body
}

// unmarshal an AccountData from a Response, encoded as a JSON []byte.
// will panic if unmarshalling fails
func unmarshal(input []byte) AccountData {
	var ret Response

	err := json.Unmarshal(input, &ret)
	if err != nil {
		panic(err)
	}

	return *ret.Data
}

type Request struct {
	Data *AccountData `json:"data"`
}

type Response struct {
	Data *AccountData `json:"data"`
}
