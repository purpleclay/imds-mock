---
icon: material/aws
---

# Instance Metadata Categories

Instance metadata is divided into categories[^1]. To retrieve instance metadata, you specify the category in the request, and the metadata is returned in the response.

The `imds-mock` offers different levels of support for each of the instance categories. Please consult this table as new features will be supported with each future release.

!!! info "Table Key"

    This is a living table and will change as new features are released. In the meantime, the following key highlights the level of support for each metadata category.

    - :material-check-all:: fully supported
    - :material-check:: partially supported, future enhancements planned
    - :material-close:: not currently supported

## Categories

The following table lists the categories of instance metadata. {==Highlighted==} text within a category denotes a dynamic placeholder.

| Category                                                              | Supported                                              |
| --------------------------------------------------------------------- | ------------------------------------------------------ |
| `ami-id`                                                              | :material-check-all:{title="fully supported"} `v0.1.0` |
| `ami-launch-index`                                                    | :material-check-all:{title="fully supported"} `v0.1.0` |
| `ami-manifest-path`                                                   | :material-close:{title="not supported"}                |
| `ancestor-ami-ids`                                                    | :material-close:{title="not supported"}                |
| `autoscaling/target-lifecycle-state`                                  | :material-close:{title="not supported"}                |
| `block-device-mapping/ami`                                            | :material-check:{title="partially supported"} `v0.1.0` |
| `block-device-mapping/ebs{==N==}`                                     | :material-check:{title="partially supported"} `v0.1.0` |
| `block-device-mapping/ephemeral{==N==}`                               | :material-close:{title="not supported"}                |
| `block-device-mapping/root`                                           | :material-check:{title="partially supported"} `v0.1.0` |
| `block-device-mapping/swap`                                           | :material-close:{title="not supported"}                |
| `elastic-gpus/associations/{==elastic-gpu-id==}`                      | :material-close:{title="not supported"}                |
| `elastic-inference/associations/{==eia-id==}`                         | :material-close:{title="not supported"}                |
| `events/maintenance/history`                                          | :material-close:{title="not supported"}                |
| `events/maintenance/scheduled`                                        | :material-close:{title="not supported"}                |
| `events/recommendations/rebalance`                                    | :material-check-all:{title="fully supported"} `v0.3.0` |
| `hostname`                                                            | :material-check-all:{title="fully supported"} `v0.1.0` |
| `iam/info`                                                            | :material-check-all:{title="fully supported"} `v0.1.0` |
| `iam/security-credentials/{==role-name==}`                            | :material-check-all:{title="fully supported"} `v0.1.0` |
| `identity-credentials/ec2/info`                                       | :material-close:{title="not supported"}                |
| `identity-credentials/ec2/security-credentials/ec2-instance`          | :material-close:{title="not supported"}                |
| `instance-action`                                                     | :material-check-all:{title="fully supported"} `v0.1.0` |
| `instance-id`                                                         | :material-check-all:{title="fully supported"} `v0.1.0` |
| `instance-life-cycle`                                                 | :material-check:{title="partially supported"} `v0.1.0` |
| `instance-type`                                                       | :material-check:{title="partially supported"} `v0.1.0` |
| `ipv6`                                                                | :material-close:{title="not supported"}                |
| `kernel-id`                                                           | :material-close:{title="not supported"}                |
| `local-hostname`                                                      | :material-check-all:{title="fully supported"} `v0.1.0` |
| `local-ipv4`                                                          | :material-check-all:{title="fully supported"} `v0.1.0` |
| `mac`                                                                 | :material-check-all:{title="fully supported"} `v0.1.0` |
| `metrics/vhostmd`                                                     | :material-close:{title="not supported"}                |
| `network/interfaces/macs/{==mac==}/device-number`                     | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/interface-id`                      | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/ipv4-associations/{==public-ip==}` | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/ipv6s`                             | :material-close:{title="not supported"}                |
| `network/interfaces/macs/{==mac==}/local-hostname`                    | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/local-ipv4s`                       | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/mac`                               | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/network-card-index`                | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/owner-id`                          | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/public-hostname`                   | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/public-ipv4s`                      | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/security-groups`                   | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/security-group-ids`                | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/subnet-id`                         | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/subnet-ipv4-cidr-block`            | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/subnet-ipv6-cidr-blocks`           | :material-close:{title="not supported"}                |
| `network/interfaces/macs/{==mac==}/vpc-id`                            | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/vpc-ipv4-cidr-block`               | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/vpc-ipv4-cidr-blocks`              | :material-check-all:{title="fully supported"} `v0.1.0` |
| `network/interfaces/macs/{==mac==}/vpc-ipv6-cidr-blocks`              | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/availability-zone`                                         | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/availability-zone-id`                                      | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/group-name`                                                | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/host-id`                                                   | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/partition-number`                                          | :material-check-all:{title="fully supported"} `v0.1.0` |
| `placement/region`                                                    | :material-check-all:{title="fully supported"} `v0.1.0` |
| `product-codes`                                                       | :material-check-all:{title="fully supported"} `v0.1.0` |
| `public-hostname`                                                     | :material-check-all:{title="fully supported"} `v0.1.0` |
| `public-ipv4`                                                         | :material-check-all:{title="fully supported"} `v0.1.0` |
| `public-keys/0/openssh-key`                                           | :material-check-all:{title="fully supported"} `v0.1.0` |
| `ramdisk-id`                                                          | :material-close:{title="not supported"}                |
| `reservation-id`                                                      | :material-check-all:{title="fully supported"} `v0.1.0` |
| `security-groups`                                                     | :material-check-all:{title="fully supported"} `v0.1.0` |
| `services/domain`                                                     | :material-check-all:{title="fully supported"} `v0.1.0` |
| `services/partition`                                                  | :material-check-all:{title="fully supported"} `v0.1.0` |
| `spot/instance-action`                                                | :material-check-all:{title="fully supported"} `v0.3.0` |
| `spot/termination-time`                                               | :material-check-all:{title="fully supported"} `v0.3.0` |
| `tags/instance`                                                       | :material-check-all:{title="fully supported"} `v0.2.0` |

## Dynamic Categories

The following table lists the categories of dynamic data.

| Category                      | Supported                               |
| ----------------------------- | --------------------------------------- |
| `fws/instance-monitoring`     | :material-close:{title="not supported"} |
| `instance-identity/document`  | :material-close:{title="not supported"} |
| `instance-identity/pkcs7`     | :material-close:{title="not supported"} |
| `instance-identity/signature` | :material-close:{title="not supported"} |

[^1]: View the official AWS documentation with regards to instance metadata categories [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html).
