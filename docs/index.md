# Instance Metadata Service Mock

The Instance Metadata Service (IMDS) stores data about an EC2 that you can use to configure or manage the running of that instance. Data is made accessible through the use of instance categories that adapt to its current state. Designing and developing software around the IMDS service poses two main challenges:

1. First and foremost, an EC2 instance is needed; increasing costs
1. There is no way to influence the IMDS service to simulate EC2 events such as spot termination

Both of which make testing difficult and unattainable.

## So why use a Mock?

The `imds-mock` attempts to solve these problems by providing a tool to accurately simulate any use case within the IMDS service, bringing testing to the forefront without additional cost.

### Features

- All mock responses accurately reflect those from the actual IMDS service
- Customisation of responses is supported through CLI flags
- Support for both IMDSv1 and IMDSv2, with strict IMDSv2 possible
- An in-built eventing system makes the simulation of spot interruption notices both easy and configurable
