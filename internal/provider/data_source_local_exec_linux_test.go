// +build linux

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestLocalExecDataSource(t *testing.T) {
	var tests = []struct {
		stdout string
		stderr string
		rc     string
		config string
	}{
		{
			"hello",
			"world",
			"0",
			`
				data "local_exec" "command" {
					command = ["sh", "-c", "echo -n hello; echo -n 1>&2 world"]
				}
			`,
		},
		{
			"",
			"",
			"127",
			`
				data "local_exec" "command" {
					command = ["sh", "-c", "exit 127"]
                    ignore_failure = true
				}
			`,
		},
		{
			"/tmp\n",
			"",
			"0",
			`
				data "local_exec" "command" {
					command = ["pwd"]
                    working_dir = "/tmp"
				}
			`,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testProviders,
				Steps: []resource.TestStep{
					{
						Config: test.config,
						Check: func(s *terraform.State) error {
							m := s.RootModule()
							i := m.Resources["data.local_exec.command"].Primary
							if got, want := i.Attributes["stdout"], test.stdout; got != want {
								return fmt.Errorf("stdout %q; want %q", got, want)
							}
							if got, want := i.Attributes["stderr"], test.stderr; got != want {
								return fmt.Errorf("stderr %q; want %q", got, want)
							}
							if got, want := i.Attributes["rc"], test.rc; got != want {
								return fmt.Errorf("rc %q; want %q", got, want)
							}
							return nil
						},
					},
				},
			})
		})
	}
}
