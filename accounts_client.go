// Provides some basic CRUD functions for interacting with the account api.
// Generally speaking, errors are divided into ClientErrors and ServerErrors, where ClientErrors represent
// status codes > 399 < 500 and ServerErrors are status codes >= 500.
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

// CreateAccount sends an account creation request to the API, returning the resulting account object.
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

// FetchAccount fetches the account with the given uuid. If the account doesn't exist the function will return (nil, nil).
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

// DeleteAccount deletes the account with the given ID and given version number. If the given resource cannot be found
// this will return a ClientError wrapping a 404 statusCode. A possibly useful improvement here would be to
// define a constant ClientError for the 404 case, because clients may want to tolerate 404's and not treat them
// at all. Didn't implement this because I didn't want to make an assumption on the usefulness of this.
func DeleteAccount(accountId string, version int64) error {
	path := fmt.Sprintf("/organisation/accounts/%v?version=%d", accountId, version)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, buildUrl(path), nil)
	if err != nil {
		return &ConnectionError{err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return &ConnectionError{err}
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ConnectionError{err}
	}

	if resp.StatusCode > 299 && resp.StatusCode < 500 {
		return &ClientError{StatusCode: resp.StatusCode, Message: string(bytes)}
	}
	if resp.StatusCode >= 500 {
		return &ServerError{StatusCode: resp.StatusCode, Message: string(bytes)}
	}

	return nil
}

// buildUrl builds a complete url from the env var BASEURL + the given path.
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

// marshal will marshal the given AccountData into a Request, as json, in a []byte.
// will panic if marshalling fails.
func marshal(input *AccountData) []byte {
	req := Payload{
		Data: input,
	}

	body, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}

	return body
}

// unmarshal will unmarshal an AccountData from a Response, encoded as a JSON []byte.
// will panic if unmarshalling fails
func unmarshal(input []byte) AccountData {
	var ret Payload

	err := json.Unmarshal(input, &ret)
	if err != nil {
		panic(err)
	}

	return *ret.Data
}

// Payload is a wrapper for the accounts api.
type Payload struct {
	Data *AccountData `json:"data"`
}
