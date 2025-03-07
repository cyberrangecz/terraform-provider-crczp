---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "crczp_sandbox_pool Resource - terraform-provider-crczp"
subcategory: ""
description: |-
  Sandbox pool
---

# crczp_sandbox_pool (Resource)

Sandbox pool

## Example Usage

```terraform
resource "crczp_sandbox_definition" "example" {
  url = "https://github.com/cyberrangecz/library-junior-hacker.git"
  rev = "master"
}

resource "crczp_sandbox_pool" "example" {
  definition = {
    id = crczp_sandbox_definition.example.id
  }
  max_size = 2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition` (Attributes) The associated sandbox definition (see [below for nested schema](#nestedatt--definition))
- `max_size` (Number) Maximum number of allocated sandbox allocation units

### Read-Only

- `created_by` (Attributes) Who created the sandbox pool (see [below for nested schema](#nestedatt--created_by))
- `hardware_usage` (Attributes) Current resource usage by all allocation units in the pool (see [below for nested schema](#nestedatt--hardware_usage))
- `id` (Number) Id of the sandbox pool
- `lock_id` (Number) Id of the associated lock
- `rev` (String) Revision of the associated Git repository used for the sandbox pool
- `rev_sha` (String) Revision hash of the associated Git repository used for the sandbox pool
- `size` (Number) Current number of allocated sandbox allocation units

<a id="nestedatt--definition"></a>
### Nested Schema for `definition`

Required:

- `id` (Number) Id of the associated sandbox definition

Read-Only:

- `created_by` (Attributes) Who created the sandbox definition (see [below for nested schema](#nestedatt--definition--created_by))
- `name` (String) Name of the sandbox definition
- `rev` (String) Revision of the Git repository of the sandbox definition
- `url` (String) Url to the Git repository of the sandbox definition

<a id="nestedatt--definition--created_by"></a>
### Nested Schema for `definition.created_by`

Read-Only:

- `family_name` (String) Family name of the user
- `full_name` (String) Full name of the user
- `given_name` (String) Given name of the user
- `id` (Number) Id of the user
- `mail` (String) Email of the user
- `sub` (String) Sub of the user as given by an OIDC provider



<a id="nestedatt--created_by"></a>
### Nested Schema for `created_by`

Read-Only:

- `family_name` (String) Family name of the user
- `full_name` (String) Full name of the user
- `given_name` (String) Given name of the user
- `id` (Number) Id of the user
- `mail` (String) Email of the user
- `sub` (String) Sub of the user as given by an OIDC provider


<a id="nestedatt--hardware_usage"></a>
### Nested Schema for `hardware_usage`

Read-Only:

- `instances` (String) The percentage of used instances relative to the cloud quota
- `network` (String) The percentage of used networks relative to the cloud quota
- `port` (String) The percentage of used ports relative to the cloud quota
- `ram` (String) The percentage of used RAM relative to the cloud quota
- `subnet` (String) The percentage of used subnets relative to the cloud quota
- `vcpu` (String) The percentage of used vCPUs relative to the cloud quota
