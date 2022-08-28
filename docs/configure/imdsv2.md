---
icon: material/shield-key-outline
---

# IMDSv2

IMDSv2 uses session-orientated requests, prohibiting access to any endpoints on the imds-mock without using a session token. A client must request a token with a maximum TTL of between one second and six hours before further requests.

It is good security practice to only support IMDSv2[^1] when launching an EC2. Enable the `--imdsv2` flag to simulate this behaviour.

## Enforce Strict IMDSv2

=== "CLI"

    ```sh
    imds-mock --imdsv2
    ```

=== "DockerHub"

    ```sh
    docker run -p 1338:1338 purpleclay/imds-mock --imdsv2
    ```

=== "GHCR"

    ```sh
    docker run -p 1338:1338 ghcr.io/purpleclay/imds-mock --imdsv2
    ```

## Using a Session Token

1. Request a session token by providing the `X-aws-ec2-metadata-token-ttl-seconds` header with a value between `1` and `21600` seconds (_six hours_):
   ```sh
   TOKEN=`curl -X PUT "http://localhost:1338/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`
   ```
1. Include the token in any subsequent requests by providing the `X-aws-ec2-metadata-token` header:
   ```sh
   curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://localhost:1338/latest/meta-data/
   ```

[^1]: The AWS Security blog post, [Add defense in depth against open firewalls, reverse proxies, and SSRF vulnerabilities with enhancements to the EC2 Instance Metadata Service](https://aws.amazon.com/blogs/security/defense-in-depth-open-firewalls-reverse-proxies-ssrf-vulnerabilities-ec2-instance-metadata-service/), details why using IMDSv2 is important to EC2 security
