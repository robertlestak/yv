# yv - YAML Validator

`yv` is a tool which enables you to validate YAML files against a defined ruleset. yv can be easily placed in CI/CD pipelines to validate manifests against enforced security policies, check for pre-validated use cases, and auto-approving PRs.

## Configuration

`yv` is configured using command line flags, `-r` for the path to the ruleset and `-f` for the path to the yamls to validate. If provided as a path to a file, only that single file will be read. However if a directory is provided, the directory and all children directories will be read. This is useful for at-scale validation.

### Example Configuration

```yaml
rules:
- type: restricted-path-values
  description: prevent the creation of resources with external gateways
  name: deny-istio-external
  rule:
    path: spec.gateways[*]
    value: "istio-system/ingressgateway-ext"
    regex: false
- type: restricted-path-values
  description: prevent the creation of service entries
  name: deny-istio-se
  rule:
    path: kind
    value: ServiceEntry
    regex: false
- type: restricted-path-values
  description: prevent the modification of system resources
  name: deny-system-resources
  rule:
    path: metadata.name
    value: ".*-system"
    regex: true
- type: globally-unique-values
  name: unique-waf-site-id
  description: ensure that the WAF site ID is unique across all resources
  rule:
    path: metadata.annotations.'cert-manager-sync.lestak.sh/incapsula-site-id'
```

## Supported Validations

### restricted-path-values

Prevent the modification of resources based on the path selector and a value. Regex boolean (default false) enables regex matching against the path value.

#### Options

**path** string: yaml path to select

**value** string: value to match

**regex** bool: regex match value (default: false)

**example**:

```yaml
- type: restricted-path-values
  description: restrict values in specific paths
  name: istio-external
  rule:
    path: spec.gateways[*]
    value: "istio-system/ingressgateway-ext"
    regex: false
```

### globally-unique-values

Prevent two resources defined by a path from having the same value. This is useful when validating WAF integrations, Istio ServiceEntries, or other globally unique resources.

#### Options

**path** string: yaml path to check globally

**example**:

```yaml
- type: globally-unique-values
  name: unique-waf-site-id
  description: ensure that the WAF site ID is unique across all resources
  rule:
    path: metadata.annotations.'cert-manager-sync.lestak.sh/incapsula-site-id'
```