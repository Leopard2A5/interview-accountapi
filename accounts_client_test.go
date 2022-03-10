package accounts

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
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
	t.Logf("'%#v', %v", resp, err)
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
