# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address    = "127.0.0.1:8080"
  reflection = true
}

client {
  address = "127.0.0.1:8080"
}

database {
  type = "sqlite"
  dsn  = "./configs/monolith.sqlite3"
}

