# Tortilla

T. is a wrapper of vault-agent

## Usage

```bash
tortilla --help
====================================================================
 It's a Wrap!

 Tortilla: A wrap-per for vault-agent to manage secrets easily.
====================================================================
Usage: tortilla <command> [arguments]
```
This will run your command with the environment variables populated from Vault secrets as defined in your `tortilla.yaml` configuration file.

## Configuration
Tortilla is configured via a `tortilla.yaml` file. Below is an example configuration:

```yaml
logLevel: info
secrets:
  - path: /secret/data/creds
transformations:
  - type: replace
    match: "CREDS_SECRETNAME"
    change: "DB_PASSWORD"
```

In this configuration:
- `logLevel` sets the logging level for Tortilla.
- `secrets` is a list of secrets to fetch from Vault. Each secret can specify:
  - `path`: The path in Vault where the secret is stored.
- `transformations` is a list of transformations to apply to the Vault Agent configuration. Each transformation can specify:
  - `type`: The type of transformation (e.g., `replace`, `prefix`).
  - `match`: The string to match for replacement.
  - `change`: The string to replace the matched string with.
  - `path`: Optionally the path to apply the prefix transformation to
Make sure to adjust the configuration according to your Vault setup and the secrets you want to manage.

## Example

Using the previos configuration file, and a test.sh like this:

```bash
#!/bin/bash
echo "DB_PASSWORD = $DB_PASSWORD"
```

You can run tortilla as follows:

```bash 
tortilla ./test.sh

Successfully generated "/var/folders/44/8pl54y0d5hs_hh94qz2yfxkc0000gn/T/vault-agent-config-1965862777.hcl" configuration file!
Warning: the generated file uses 'token_file' authentication method, which is not suitable for production environments.
==> Vault Agent started! Log data will stream in below:

==> Vault Agent configuration:

           Api Address 1: http://bufconn
                     Cgo: disabled
               Log Level: error
                 Version: Vault v1.20.4, built 2025-09-23T13:22:38Z
             Version Sha: 55bd8f18c6c84aa89fdede4850a622c57f03bd7e

DB_PASSWORD = my-long-password
```