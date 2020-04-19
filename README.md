# kms-vault

kms-vault allows you to manage configuration files that partially contain confidential information in a repository using Cloud KMS.
You can apply [Ansible vault best practices](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html#variables-and-vaults) and encryption key sharing for team development using Cloud KMS. The commands are made almost compatible for `ansible-vault`, so the migration is easy.

# Install

```shell script
curl -sfL https://raw.githubusercontent.com/juntaki/kms-vault/master/install.sh | sh
```

# Basic usage

Follow the steps in [Cloud KMS Documentation](https://cloud.google.com/kms/docs/creating-keys) to create a key ring and a key.
Write the string you want to encrypt into a file to encrypt it.

```shell script
> echo secret! > secret-file
> cat secret-file
secret!

> # Encryption
> kms-vault encrypt --project <your-project-name> --location <location> --keyring <your-keyring> --key <your-key> secret-file
> cat secret-file
$VAULT;0.1.0;CLOUD_KMS
CiQAX7XZ1ruzQYWP6my24jdc6BHt+tMifMzQVEOi1QROl4fzyAMSMQDKS25uyx7kqF8v/VcqwV3n2mhzl9Wm13voO5dgb0EeJZk489GAV9RoWWuinHeJxhE=

> # Decription
> kms-vault decrypt --project <your-project-name> --location <location> --keyring <your-keyring> --key <your-key> secret-file
> cat secret-file
secret!
```

# Save the repository secret

## kms-vault configration

Project, location, keyring, and key can be saved as configuration file for the repository.
Put the file in the project root.

```shell script
> # Write your config to .kms-vault.yaml
> kms-vault config --project <your-project-name> --location <location> --keyring <your-keyring> --key <your-key> -w
> cat .kms-vault.yaml 
project: your-project-name
location: location
keyring: your-keyring
key: your-key
```

## Ansible best practice

Only the secret information is extracted and saved to vault.yaml in YAML format.
vars.yaml is a template configuration file that does not contain any secret information.

```
├── vars.yaml
└── vault.yaml
```

vault.yaml gives a name to each of the secrets. (The extension must be yaml or yml in this case.)
To see the vault.yaml in an encrypted state, use the `view` command. You can use the `--yaml` option to check the values used for embedding a template. The file name is lower-cased without the extension, and the namespace is automatically split.

```shell script
> kms-vault view vault.yaml
secret: secret-text

> kms-vault view --yaml vault.yaml
vault:
  secret: secret-text
```

vars.yaml is Go's text/template format, and you can get the output with the `fill` command.

```shell script
> cat vars.yaml 
plain_text: not-secret-text
secret: {{.vault.secret}}
> kms-vault fill --template vars.yaml vault.yaml 
plain_text: not-secret-text
secret: secret-text
```
