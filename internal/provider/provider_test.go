package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testProviders = map[string]*schema.Provider{
	"local": New(),
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func checkFileDeleted(shouldNotExistFile string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if _, err := os.Stat(shouldNotExistFile); os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("file %s was not deleted", shouldNotExistFile)
	}
}
