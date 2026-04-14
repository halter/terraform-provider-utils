# terraform-provider-utils

A Terraform provider that exposes utility functions for CIDR/IP network operations.

## Functions

### `cidrcontains(containing_prefix, contained_prefix)`

Returns `true` if the containing CIDR prefix fully encompasses the contained CIDR prefix.

```hcl
provider::utils::cidrcontains("10.0.0.0/8", "10.1.0.0/16")    # true
provider::utils::cidrcontains("10.0.0.0/8", "192.168.0.0/16")  # false
provider::utils::cidrcontains("10.0.0.0/8", "10.1.2.3/32")     # true
```

### `cidroverlaps(prefix_a, prefix_b)`

Returns `true` if two CIDR prefixes have any addresses in common.

```hcl
provider::utils::cidroverlaps("10.0.0.0/8", "10.1.0.0/16")     # true
provider::utils::cidroverlaps("10.0.0.0/8", "192.168.0.0/16")  # false
```

### `cidrnoverlap(prefixes)`

Returns `true` if no pair of CIDR prefixes in the list overlaps with each other.

```hcl
provider::utils::cidrnoverlap(["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"])  # true
provider::utils::cidrnoverlap(["10.0.0.0/8", "10.1.0.0/16", "192.168.0.0/16"])    # false
```

All functions:
- Require CIDR notation for all arguments (use `/32` or `/128` for single hosts)
- Support both IPv4 and IPv6
- Return an error if address families are mixed

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.24

## Building the Provider

```shell
go install
```

## Developing the Provider

To generate documentation:

```shell
make generate
```

To run tests:

```shell
make testacc
```
