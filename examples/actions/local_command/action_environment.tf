resource "terraform_data" "notify" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.notify]
    }
  }
}

action "local_command" "notify" {
  config {
    command   = "bash"
    arguments = ["scripts/notify.sh"]
    environment = {
      APP_NAME    = var.app_name
      ENVIRONMENT = var.environment
      LOG_LEVEL   = var.log_level
    }
  }
}
