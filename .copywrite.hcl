# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0

schema_version = 1

project {
  license = "MPL-2.0"

  copyright_year   = 2025
  copyright_holder = "Kasefuchs"

  header_ignore = [
    # --- IDE / Editor ---
    ".idea/**",
    ".vscode/**",

    # --- Golang ---
    # Protobuf artifacts
    "**/*.pb.go",

    # --- Third party ---
    "third_party/**"
  ]
}
