---
icon: material/timer-cog-outline
status: new
---

# Spot Instance

The imds-mock can simulate a spot instance. A spot instance is more cost-effective than an on-demand instance but is at the mercy of spot interruptions. An interruption notice will typically be issued by Amazon EC2 two minutes before it stops or terminates your spot instance, with no warning issued before hibernation.

## Switching Instance Type

Set the `--spot` flag to enable spot simulation within the imds-mock. The `spot/instance-action` metadata category will immediately be available and will return an interruption notice to terminate the instance.

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

You have complete control over the spot interruption notice raised by the imds-mock. Set the `--spot-action` flag, specifying the interruption type (`terminate`, `stop` or `hibernate`) and an initial delay to raise a spot interruption notice.

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

    A hibernate interruption notice does not provide a two-minute warning and is effective immediately. It, therefore, should not be accessible through the `spot/instance-action` metadata category. However, as the mock will remain running, this category will be available and contain details of the hibernation interruption.
