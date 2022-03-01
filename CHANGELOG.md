## 2.2.0 (March 1, 2022)

### New features

* Added optional argument `sensitive` to data source `local_file`. When used attributes `content` and `content_base64` will be empty,
  and instead `sensitive_content` and `sensitive_content_base64` will be populated. Those are defined in the schema as
  `Sensitive: true`, so Terraform will omit them from outputs ([#101](https://github.com/hashicorp/terraform-provider-local/pull/101)).


## 2.1.0 (February 19, 2021)

### Notes

* Binary releases of this provider now include the darwin-arm64 platform. This version contains no further changes.


## 2.0.0 (October 14, 2020)

### Notes

* Binary releases of this provider now include the linux-arm64 platform.

### Breaking changes

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. ([#42](https://github.com/terraform-providers/terraform-provider-local/issues/42))

### New features

* Add `source` attribute to `local_file` resource ([#44](https://github.com/terraform-providers/terraform-provider-local/issues/44))


## 1.4.0 (September 30, 2019)

### Notes

* The provider has switched to the standalone TF SDK, there should be no noticeable impact on compatibility. ([#32](https://github.com/terraform-providers/terraform-provider-local/issues/32))

### New features

* `local_file`: allow for configurable permissions ([#30](https://github.com/terraform-providers/terraform-provider-local/issues/30))


## 1.3.0 (June 26, 2019)

### New features

* Add support for base64 encoded content ([#29](https://github.com/terraform-providers/terraform-provider-local/issues/29))


## 1.2.2 (May 01, 2019)

### Notes

* This releases includes another Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.


## 1.2.1 (April 11, 2019)

### Notes

* This releases includes only a Terraform SDK upgrade intended to align with that being used for other providers as we prepare for the Core v0.12.0 release. It should have no significant changes in behavior for this provider.


## 1.2.0 (March 20, 2019)

### New features

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.
* `local_file` resource has optional `sensitive_content` attribute, which can be used instead of `content` in situations where the content contains sensitive information that should not be displayed in a rendered diff. ([#9](https://github.com/terraform-providers/terraform-provider-local/issues/9))


## 1.1.0 (January 04, 2018)

### New features

* `local_file` data source, for reading files in a way that participates in Terraform's dependency graph, which allows reading of files that are created dynamically during `terraform apply`. ([#6](https://github.com/terraform-providers/terraform-provider-local/issues/6))


## 1.0.0 (September 15, 2017)

* No changes from 0.1.0; just adjusting to [the new version numbering scheme](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning/).


## 0.1.0 (June 21, 2017)

### Notes

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
