# Copyright IBM Corp. 2017, 2026
# SPDX-License-Identifier: MPL-2.0

variable "bash_path" {
  type    = string
  default = "bash"
}

variable "var1" {
  type = string
}

variable "var2" {
  type = string
}

variable "working_directory" {
  type    = string
  default = null
}

variable "scripts_folder_path" {
  type    = string
  default = null
}
