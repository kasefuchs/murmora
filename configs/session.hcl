# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address    = "127.0.0.1:8083"
  reflection = true
}

database {
  type = "sqlite"
  dsn  = "./configs/session.sqlite3"
}

user_service {
  address = "127.0.0.1:8081"
}

token_service {
  address = "127.0.0.1:8082"
}
