---
icon: material/timer-cog-outline
status: new
---

# Spot Instance

The imds-mock can simulate a spot instance. A spot instance is more cost-effective than an on-demand instance, but is at the mercy of spot interruptions. Typically an interruption notice will be issued by Amazon EC2 two minutes before it either stops or terminates your spot instance. No warning is issued prior to hibernation.

## Switching Instance Type

To enable simulation of a spot instance the `--spot` flag should be used. The `spot/instance-action` metadata category will immediately be available and will return an interruption notice to terminate the instance.

=== "CLI"

    ```sh
    imds-mock --spot
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock --spot
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock --spot
    ```

## Configure Interruption Notice

You have full control over the type of interruption notice that is raised within the imds-mock. The type (`terminate`, `stop` and `hibernate`) and initial delay for raising the interruption notice can be configured through the `--spot-action` flag.

=== "CLI"

    ```sh
    imds-mock --spot --spot-action stop=10s
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock --spot --spot-action stop=10s
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock --spot --spot-action stop=10s
    ```

!!! info "Handling hibernation a little differently"

    A hibernate interruption notice does not provide a two minute warning and is effective immediately. It therefore should not be accessible through the `spot/instance-action` metadata category. However, as the mock will remain running, this category will be available and contain details of when the hibernation was initiated.
