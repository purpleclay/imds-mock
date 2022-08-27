---
icon: material/tag-outline
status: new
---

# Instance Tags

EC2 instance tags can be exposed through the AWS Instance Metadata Service through the `tags/instance` instance category. The imds-mock exposes a default `Name=imds-mock-ec2` tag to simulate the enablement of this feature.

## Custom Tags

If you wish to override the default instance tags exposed by the imds-mock the `--instance-tags` flag accepts a list of `key=value` pairs.

=== "CLI"

    ```sh
    imds-mock --instance-tags Name=Test,Environment=Dev
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock --instance-tags Name=Test,Environment=Dev
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock --instance-tags Name=Test,Environment=Dev
    ```

### Querying a Tag

Any custom tag can be retrieved using the root `tags/instance` metadata category. To therefore retrieve the custom `Environment` tag, the following request would be used:

```sh
curl http://localhost:1338/latest/meta-data/tags/instance/Environment
```

## Excluding Instance Tags

By default EC2 instance tags are omitted from the AWS Instance Metadata Service. To simulate this in the imds-mock, tags can be explicitly excluded through the `--exclude-instance-tags` flag.

=== "CLI"

    ```sh
    imds-mock --exclude-instance-tags
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock --exclude-instance-tags
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock --exclude-instance-tags
    ```
