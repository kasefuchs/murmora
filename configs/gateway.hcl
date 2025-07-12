# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

server {
  address = "127.0.0.1:8079"
  prefix  = "/gateway/v1"
}

authentication_service {
  address = "127.0.0.1:8084"
}
