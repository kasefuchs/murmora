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

    # --- Node ---
    # Package stores
    ".pnpm-store/**",
    "node_modules/**",

    # Lock files
    "pnpm-lock.yaml",

    # --- Third party ---
    "third_party/**"
  ]
}
