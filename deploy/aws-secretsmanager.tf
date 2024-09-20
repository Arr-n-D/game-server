resource "random_password" "db_password"{
  length           = 18
  special          = true
  override_special = "_!%^"
}

resource "aws_secretsmanager_secret" "db_secrets" {
  name = "db"
}

resource "aws_secretsmanager_secret_version" "db_secrets" {
  secret_id = aws_secretsmanager_secret.db_secrets.id
  secret_string = random_password.db_password.result
}

