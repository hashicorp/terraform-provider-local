package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"local_file":           resourceLocalFile(),
			"local_sensitive_file": resourceLocalSensitiveFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"local_file":           dataSourceLocalFile(),
			"local_sensitive_file": dataSourceLocalSensitiveFile(),
		},
	}
}
