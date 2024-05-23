# Vault Raft Backup

[![Build Status](https://zuul.otc-service.com/api/tenant/eco/badge?project=opentelekomcloud-infra/vault-raft-backup&pipeline=gate)](https://zuul.otc-service.com/t/eco/builds?project=opentelekomcloud-infra%2Fvault-raft-backup&pipeline=gate)
[![GitHub release](https://img.shields.io/github/v/release/opentelekomcloud-infra/vault-raft-backup.svg)](https://github.com/opentelekomcloud-infra/vault-raft-backup/releases/latest)
[![GitHub license](https://img.shields.io/github/license/opentelekomcloud-infra/vault-raft-backup.svg)](https://github.com/opentelekomcloud-infra/vault-raft-backup/blob/master/LICENSE)

`vault-raft-backup` is a command line utility designed to manage HashiCorp Vault backups using [raft snapshots](https://learn.hashicorp.com/tutorials/vault/sop-backup) and transfer these backups to Object Storage Service (OBS).

Enabling version control in the OBS bucket is crucial as it not only saves each specified backup but also maintains a copy under the "latest" tag, ensuring easy retrieval of the most recent backup.

## Prerequisites

Before using `vault-raft-backup`, ensure the following conditions are met:

- Proper authentication and authorization credentials for an OTC account to perform list, read, and write operations in an OBS bucket.
- Access to an OBS bucket with versioning enabled to store snapshots.
- Connectivity to a functional Vault server cluster that utilizes Raft integrated storage.
- Sufficient permissions on the Vault server cluster to execute Raft snapshots, typically provided via an AppRole with the necessary policies.

Here is an example of a Vault policy needed to perform Raft snapshots:

```hcl
path "sys/storage/raft/snapshot" {
  capabilities = ["read"]
}

```
## Install

You can install vault-backup using the following go command:

    go install github.com/opentelekomcloud-infra/vault-raft-backup@latest

Or use docker image:

    docker pull quay.io/opentelekomcloud/vault-raft-backup:release_0.0.1


## List of flags

You can use flags, environment variables or config file to set the same values.

| Variable        | Flag              | Description                       | Required | Default                                |
|-----------------|-------------------|-----------------------------------|----------|----------------------------------------|
| OBS_BUCKET_NAME | --obs-bucket-name | OBS bucket name                   | true     |                                        |
| OBS_OBJECT_NAME | --obs-object-name | OBS object name                   | false    | vault-raft-backup.snap                 |
| OS_ACCESS_KEY   | --os-access-key   | OTC Access Key for authentication | true     |                                        |
| OS_SECRET_KEY   | --os-secret-key   | OTC Secret Key for authentication | true     |                                        |
| OS_AUTH_URL     | --os-auth-url     | OTC Authentication URL            | false    | https://iam.eu-de.otc.t-systems.com/v3 |
| OS_DOMAIN_NAME  | --os-domain-name  | OTC Domain name                   | false    | eu-de                                  |
| OS_PROJECT_NAME | --os-project-name | OTC Project name                  | true     |                                        |
| VAULT_ADDRESS   | --vault-address   | Vault address                     | false    | https://127.0.0.1:8200                 |
| VAULT_ROLE_ID   | --vault-role-id   | Vault AppRole role ID             | true     |                                        |
| VAULT_SECRET_ID | --vault-secret-id | Vault AppRole secret ID           | true     |                                        |
| VAULT_TIMEOUT   | --vault-timeout   | Vault Client timeout              | false    | 60s                                    |
| CONFIG          | --config          | Config file                       | false    | `$HOME/.vault-raft-backup.yaml`        |

For more details use `vault-raft-backup backup --help`

## Example usage

```
docker run --rm -it vault-raft-backup backup --vault-role-id xxx --vault-secret-id xxx --obs-bucket-name "test-bucket" --os-access-key xxx  --os-secret-key xxx --os-domain-name xxx
```

```
./vault-raft-backup backup --vault-role-id xxx --vault-secret-id xxx --obs-bucket-name "test-bucket" --os-access-key xxx  --os-secret-key xxx --os-domain-name xxx
```

You can also use a config file to set the same values instead of passing them as flags or environment variables. The default config file location is $HOME/.vault-raft-backup.yaml, but you can specify a different file using the --config flag.

Example of a config file (vault-raft-backup.yaml):

```yaml
obs_bucket_name: "test-bucket"
obs_object_name: "vault-raft-backup.snap"
os_access_key: "your-access-key"
os_secret_key: "your-secret-key"
os_auth_url: "https://iam.eu-de.otc.t-systems.com/v3"
os_domain_name: "your-domain-name"
os_project_name: "your-project-name"
vault_address: "https://127.0.0.1:8200"
vault_role_id: "your-vault-role-id"
vault_secret_id: "your-vault-secret-id"
vault_timeout: "60s"
```