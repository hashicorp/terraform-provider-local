## 2.2.0 (March 10, 2022)

NOTES:

* resource/local_file: Argument `sensitive_content` is `Deprecated`. For creating or accessing files containing sensitive data,
  please use the new resource and data source `local_sensitive_file`.
  Both are identical to their `local_file` counterparts, but `content` and `content_base64` attributes are marked as _sensitive_.

FEATURES:

* **New Data Source:** `local_sensitive_file` ([#101](https://github.com/hashicorp/terraform-provider-local/pull/101) and [#106](https://github.com/hashicorp/terraform-provider-local/pull/106))
* **New Resource:** `local_sensitive_file** ([#106](https://github.com/hashicorp/terraform-provider-local/pull/106))

## 2.1.0 (February 19, 2021)

NOTES:

* Binary releases of this provider now include the` darwin-arm64` platform.
* This version contains no further changes.

## 2.0.0 (October 14, 2020)

NOTES:

* Binary releases of this provider now include the `linux-arm64` platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11.
  This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version.
  ([#42](https://github.com/terraform-providers/terraform-provider-local/issues/42))

FEATURES:

* resource/local_file: Added `source` attribute as alternative way to provide content
  for the `local_file` resource.
  ([#44](https://github.com/terraform-providers/terraform-provider-local/issues/44))

## 1.4.0 (September 30, 2019)

NOTES:

* The provider has switched to the standalone TF SDK, there should be no noticeable impact on compatibility.
  ([#32](https://github.com/terraform-providers/terraform-provider-local/issues/32))

FEATURES:

* resource/local_file: Added support for configurable permissions
  ([#30](https://github.com/terraform-providers/terraform-provider-local/issues/30))

## 1.3.0 (June 26, 2019)

FEATURES:

* resource/local_file: Added support for base64 encoded content
  ([#29](https://github.com/terraform-providers/terraform-provider-local/issues/29))
* data-source/local_file: Added support for base64 encoded content
  ([#29](https://github.com/terraform-providers/terraform-provider-local/issues/29))

## 1.2.2 (May 01, 2019)

NOTES:

* This releases includes another Terraform SDK upgrade intended to align with that being used for other providers
  as we prepare for the Core `v0.12.0` release. It should have no significant changes in behavior for this provider.

## 1.2.1 (April 11, 2019)

NOTES:

* This releases includes only a Terraform SDK upgrade intended to align with that being used for other providers
  as we prepare for the Core `v0.12.0` release. It should have no significant changes in behavior for this provider.

## 1.2.0 (March 20, 2019)

FEATURES:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.
* resource/local_file: added optional `sensitive_content` attribute, which can be used instead of `content`
  in situations where the content contains sensitive information that should not be displayed in a rendered diff.
  ([#9](https://github.com/terraform-providers/terraform-provider-local/issues/9))

## 1.1.0 (January 04, 2018)

FEATURES:

* data-source/local_file: Added for reading files in a way that participates in Terraform's dependency graph,
  which allows reading of files that are created dynamically during `terraform apply`.
  ([#6](https://github.com/terraform-providers/terraform-provider-local/issues/6))

## 1.0.0 (September 15, 2017)

NOTES:

* No changes from 0.1.0; just adjusting to
  [the new version numbering scheme](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning/).

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8.
  Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
