package accounts

import (
	"errors"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
)

func TestCreateClientError(t *testing.T) {
	input := AccountData{}
	var clientError *ClientError

	_, err := CreateAccount(&input)
	if !errors.As(err, &clientError) {
		t.Fatalf("expected a clientError but got %T", err)
	}
}

func TestCreateHappyCase(t *testing.T) {
	input := newAccount()

	resp, err := CreateAccount(&input)
	if err != nil {
		t.Fatalf("expected err to be nil, but got %v", err)
	}
	if resp.Version == nil {
		t.Fatalf("expected response to be enriched by API")
	}
}

func TestFetchAccountHappyCase(t *testing.T) {
	input := newAccount()

	account, err := CreateAccount(&input)
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	resp, err := FetchAccount(account.ID)
	if err != nil {
		t.Fatalf("failed to fetch account: %v", err)
	}

	if diff := deep.Equal(resp, account); diff != nil {
		t.Error(diff)
	}
}

func TestFetchAccountInvalidId(t *testing.T) {
	_, err := FetchAccount("not-a-uuid")

	var clientError *ClientError
	if errors.As(err, &clientError) {
		t.Fatalf("expected ClientError, but got %T", err)
	}
}

func TestFetchAccountNotFound(t *testing.T) {
	resp, err := FetchAccount(uuid.NewString())
	if resp != nil || err != nil {
		t.Fatalf("expected nil, nil, but got: %v, %v", resp, err)
	}
}

func TestBuildUrlShouldCorrectlyConcat(t *testing.T) {
	const expectation = "http://localhost:8080/v1/foo"

	os.Setenv("BASEURL", "http://localhost:8080/v1")
	if result := buildUrl("foo"); result != expectation {
		t.Fatalf("expected '%v' but got '%v'", expectation, result)
	}
	if result := buildUrl("/foo"); result != expectation {
		t.Fatalf("expected '%v' but got '%v'", expectation, result)
	}

	os.Setenv("BASEURL", "http://localhost:8080/v1/")
	if result := buildUrl("foo"); result != expectation {
		t.Fatalf("expected '%v' but got '%v'", expectation, result)
	}
	if result := buildUrl("/foo"); result != expectation {
		t.Fatalf("expected '%v' but got '%v'", expectation, result)
	}
}

func newAccount() AccountData {
	country := "GB"
	return AccountData{
		Type:           "accounts",
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Attributes: &AccountAttributes{
			Country:      &country,
			BaseCurrency: "GBP",
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			Bic:          "NWBKGB22",
			Name:         []string{"Hans Moleman"},
		},
	}
}
