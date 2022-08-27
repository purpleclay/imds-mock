---
icon: material/rocket-launch-outline
---

# On-Demand Instance

The imds-mock simulates an on-demand instance by default. Once the mock has started, all supported instance categories[^1] will be available for querying via both IMDSv1 and IMDSv2.

=== "CLI"

    ```sh
    imds-mock
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock
    ```

[^1]: A list of currently supported instance categories can be found [here](../reference/instance-metadata.md)
