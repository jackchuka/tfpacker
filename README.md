<div align="center">

# tfpacker

**Split a sprawling `main.tf` into files that make sense.**

Describe where each Terraform block belongs in a short YAML file,<br>
and tfpacker writes them out as clean, grouped `.tf` files.

[![Test](https://github.com/jackchuka/tfpacker/actions/workflows/test.yml/badge.svg)](https://github.com/jackchuka/tfpacker/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/jackchuka/tfpacker?sort=semver)](https://github.com/jackchuka/tfpacker/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

<img src="assets/demo.gif" alt="tfpacker splitting a single main.tf into storage, IAM, database and network files" width="800">

[Installation](#installation) · [Quick Start](#quick-start) · [Usage](#usage) · [Configuration](#configuration) · [Example](#example)

</div>

---

Terraform doesn't care which file a block lives in, so a growing module tends to collect everything in one `main.tf`, or in files split along lines nobody remembers. **tfpacker** reads every `.tf` file in a directory and regroups the blocks by type, provider resource type, or name, following rules you write once and can rerun on any module.

Blocks are copied byte for byte, so formatting and comments are preserved: comments inside a block, on the lines directly above it and after its closing brace move with it. Your source files are never modified.

## Features

- **Rule-based grouping** — match blocks by block type, resource type (`aws_s3_bucket`), name prefix, or name regex
- **Cross-type rules** — `ignore_type: true` groups a resource, variable and output into one file purely by name
- **Sensible defaults** — anything no rule claims goes to `variables.tf`, `outputs.tf`, `resource_<type>.tf`, and so on
- **All block types** — `resource`, `data`, `module`, `variable`, `output`, `locals`, `provider`, plus `terraform`, `moved`, `import`, `check` and others
- **Dry run** — `--dry-run` lists the files and block counts without writing anything
- **Non-destructive** — writes to a separate output directory and leaves the originals in place
- **Single binary** — for Linux, macOS and Windows

## Installation

### Homebrew

```bash
brew install jackchuka/tap/tfpacker
```

### Go

```bash
go install github.com/jackchuka/tfpacker@latest
```

### Binary

Download a prebuilt binary from the [releases page](https://github.com/jackchuka/tfpacker/releases).

### From Source

```bash
git clone https://github.com/jackchuka/tfpacker.git
cd tfpacker
make build # writes build/tfpacker
```

## Quick Start

1. **Preview** how the blocks in a module would be grouped with the default rules:

   ```bash
   tfpacker ./infra --dry-run
   ```

2. **Add rules** in `tfpacker.config.yaml` in the directory you run tfpacker from (see [Configuration](#configuration)):

   ```yaml
   rules:
     - match_type: "resource"
       sub_type: "aws_s3_bucket"
       output_file: "storage.tf"
   ```

3. **Write** the grouped files to a new directory, review them, then swap them in for the originals:

   ```bash
   tfpacker ./infra --output ./infra-packed
   ```

## Usage

```bash
tfpacker [directory] [flags]
```

`directory` defaults to the current directory. Only `.tf` files directly inside it are read; subdirectories are not scanned.

| Flag                | Description                                                                 | Default                |
| ------------------- | --------------------------------------------------------------------------- | ---------------------- |
| `--config <path>`   | Config file; built-in defaults are used if it doesn't exist                 | `tfpacker.config.yaml` |
| `--output <dir>`    | Directory to write the grouped files to                                     | `output`               |
| `--dry-run`         | Print the files that would be written, without writing them                 | `false`                |
| `--exclude <globs>` | Comma-separated file name patterns to skip, e.g. `backend.tf,*_override.tf` | (none)                 |
| `--verbose`         | Enable debug logging                                                        | `false`                |

`--config` and `--output` are resolved relative to the current working directory, not to `directory`. Existing files in the output directory with the same name are overwritten.

Other commands:

| Command               | Description                          |
| --------------------- | ------------------------------------ |
| `tfpacker version`    | Print version, commit and build date |
| `tfpacker completion` | Generate a shell completion script   |

## Configuration

`tfpacker.config.yaml` holds an ordered list of rules. A complete example is in [`example.tfpacker.config.yaml`](example.tfpacker.config.yaml).

```yaml
rules:
  # By resource type
  - match_type: "resource"
    sub_type: "aws_s3_bucket"
    output_file: "storage.tf"

  # By name prefix within one block type
  - match_type: "variable"
    name_prefix: "db_"
    output_file: "database.tf"

  # By name regex, across all block types
  - name_regex: "^db_.*|.*_database$"
    ignore_type: true
    output_file: "all_database.tf"

  # By name regex within one block type
  - match_type: "resource"
    name_regex: "^api_.*"
    output_file: "api_resources.tf"
```

### Rules

Rules are evaluated top to bottom and the first match wins. Every field a rule sets must match.

| Field         | Description                                                                                                           |
| ------------- | --------------------------------------------------------------------------------------------------------------------- |
| `match_type`  | Block type: `resource`, `data`, `module`, `variable`, `output`, `locals`, `provider`, `terraform`, …                  |
| `sub_type`    | Resource or data source type, e.g. `aws_s3_bucket`. Only `resource` and `data` blocks have one                        |
| `name_prefix` | Block name starts with this string, e.g. `prod_` matches `prod_bucket`                                                |
| `name_regex`  | Block name matches this [Go regular expression](https://pkg.go.dev/regexp/syntax) (unanchored unless you add `^`/`$`) |
| `ignore_type` | `true` matches any block type, so `match_type` is not needed                                                          |
| `output_file` | File to write matching blocks to, relative to `--output`. Subdirectories are created as needed                        |

A block's name is its last label: `bucket1` in `resource "aws_s3_bucket" "bucket1"`, `vpc` in `module "vpc"`. `locals` and `terraform` blocks have no name, so `name_prefix` and `name_regex` never match them.

### Default grouping

Blocks that no rule matches are grouped like this:

| Block         | File                                                   |
| ------------- | ------------------------------------------------------ |
| `resource`    | `resource_<type>.tf`, e.g. `resource_aws_s3_bucket.tf` |
| `data`        | `data_<type>.tf`, e.g. `data_aws_region.tf`            |
| `module`      | `modules.tf`                                           |
| `variable`    | `variables.tf`                                         |
| `output`      | `outputs.tf`                                           |
| `locals`      | `locals.tf`                                            |
| `provider`    | `providers.tf`                                         |
| anything else | `<block type>.tf`, e.g. `terraform.tf`, `moved.tf`     |

Within each output file, blocks keep the order they had in the source files (files are read in alphabetical order).

## Example

Given this `main.tf`:

```hcl
resource "aws_s3_bucket" "bucket1" { ... }
resource "aws_iam_role" "role1" { ... }
variable "region" { ... }
variable "db_password" { ... }
module "vpc" { ... }
```

With no config file, tfpacker writes:

```
output/
├── resource_aws_s3_bucket.tf   # bucket1
├── resource_aws_iam_role.tf    # role1
├── variables.tf                # region, db_password
└── modules.tf                  # vpc
```

With the rules from [Configuration](#configuration) plus one sending `aws_iam_role` to `iam.tf`:

```
output/
├── storage.tf                  # bucket1
├── iam.tf                      # role1
├── variables.tf                # region
├── database.tf                 # db_password (variable with the db_ prefix)
└── modules.tf                  # vpc
```

More input and expected output pairs live in [`test/`](test).

## Limitations

- Comments separated from a block by a blank line, such as a file header, are not copied to the output.
- Only `.tf` files are read; `.tf.json` files are ignored.

## Development

```bash
make test        # go test ./...
make build       # build/tfpacker
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT](LICENSE)
