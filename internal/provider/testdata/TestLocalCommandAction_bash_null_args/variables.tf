variable "bash_path" {
  type    = string
  default = "bash"
}

variable "stdin" {
  type    = string
  default = null
}

variable "arguments" {
  type    = list(string)
  default = []
}

variable "working_directory" {
  type    = string
  default = null
}

variable "scripts_folder_path" {
  type    = string
  default = null
}
