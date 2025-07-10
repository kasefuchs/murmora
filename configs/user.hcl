# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address    = "127.0.0.1:8081"
  reflection = true
}

database {
  type = "sqlite"
  dsn  = "./configs/user.sqlite3"
}
