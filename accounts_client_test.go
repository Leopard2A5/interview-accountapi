package accounts

import (
	"errors"
	"os"
	"testing"

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
	country := "GB"
	input := AccountData{
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

	resp, err := CreateAccount(&input)
	if err != nil {
		t.Fatalf("expected err to be nil, but got %v", err)
	}
	if resp.Version == nil {
		t.Fatalf("expected response to be enriched by API")
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
