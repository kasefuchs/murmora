# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address    = "127.0.0.1:8084"
  reflection = true
}

database {
  type = "sqlite"
  dsn  = "./configs/authentication.sqlite3"
}

user_service {
  address = "127.0.0.1:8081"
}

session_service {
  address = "127.0.0.1:8083"
}
