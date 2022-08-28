---
icon: material/console
---

# Command Line

Mocks the Amazon Instance Metadata Service (IMDS) for EC2.

## Usage

```text
imds-mock [flags]
imds-mock [command]
```

## Flags

```text
    --exclude-instance-tags          exclude access to instance tags associated with the instance
-h, --help                           help for imds-mock
    --imdsv2                         enforce IMDSv2 requiring all requests to contain a valid metadata token
    --instance-tags stringToString   a list of instance tags (key pairs) to expose as metadata (default [Name=imds-mock-ec2])
    --port int                       the port to be used at startup (default 1338)
    --pretty                         if instance categories should return pretty printed JSON
    --spot                           enable simulation of a spot instance and interruption notice
    --spot-action stringToString     configure the type and delay of the spot interruption notice (default terminate=0s)
```

## Commands

```text
completion  Generate a completion script for your target shell
help        Help about any command
version     Prints the build time version information
```
