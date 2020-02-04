package local

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-provider-local/local/test"
)

var testProviders = map[string]terraform.ResourceProvider{
	"local": Provider(),
}

func TestHello(t *testing.T) {
	t.Log(test.Hello())
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
