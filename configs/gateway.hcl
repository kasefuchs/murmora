# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address = "127.0.0.1:8079"
  prefix  = "/gateway/v1"
}

user_service {
  address = "127.0.0.1:8081"
}

token_service {
  address = "127.0.0.1:8082"
}

session_service {
  address = "127.0.0.1:8083"
}

authentication_service {
  address = "127.0.0.1:8084"
}
