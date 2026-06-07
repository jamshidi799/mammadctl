# mammadctl

`mammadctl` is a small CLI for Kubernetes helpers and Sotoon API access management.

## Installation

### Prerequisites

- Go `1.26.0` or newer, as defined in `go.mod`
- Access to a Kubernetes cluster for Kubernetes-related commands
- Sotoon API credentials for Sotoon-related commands

### Build from source

```bash
git clone <repository-url>
cd mammadctl
go build -o mammadctl .
```

You can then run it from the project directory:

```bash
./mammadctl --help
```

### Install with Go

From the project directory:

```bash
go install .
```

Make sure your Go binary directory is in your `PATH`. It is usually one of these:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

After that, run:

```bash
mammadctl --help
```

## Environment variables

### Required for Sotoon commands

Export this before using any `mammadctl sotoon ...` command:

```bash
export BEPA_TOKEN="<your-bepa-token>"
```

### Required for `sotoon access s3`

Export these before using `mammadctl sotoon access s3 ...`:

```bash
export S3_ACCESS_KEY="<your-s3-access-key>"
export S3_SECRET_KEY="<your-s3-secret-key>"
```

### Optional for Kubernetes commands

Kubernetes commands use `KUBECONFIG` if it is set. Otherwise, they use `~/.kube/config`.

```bash
export KUBECONFIG="/path/to/kubeconfig"
```

Commands that use Kubernetes config:

- `mammadctl decode secret`
- `mammadctl convert`

## Commands

```text
mammadctl
├── decode
│   └── secret
├── convert
└── sotoon
    ├── user
    │   ├── list
    │   └── access
    │       ├── cdn
    │       ├── k8s
    │       └── s3
    └── access
        ├── cdn
        ├── k8s
        └── s3
```

## Usage

### `mammadctl decode secret`

Finds a Kubernetes secret and writes it to `./secret.yaml` with decoded values in `stringData` mode.

It's useful when you want to edit a `secret` or `sealedSecret`.

```bash
mammadctl decode secret -n <namespace> <secret-name>
```

Options:

- `-n, --namespace`: Kubernetes namespace. Default: `default`

Example:

```bash
mammadctl decode secret -n production app-secret
```

Output file:

```text
./secret.yaml
```

### `mammadctl convert`

Reads a Kubernetes Service from the current cluster and writes a melange-compatible(Bazaar operator for mirroring services between clusters) formated Service manifest to `<namespace>-<service-name>.yaml`.

```bash
mammadctl convert -n <namespace> -c <cluster> <service-name>
```

Options:

- `-n, --namespace`: Kubernetes namespace. Default: `default`
- `-c, --cluster`: Cluster name to put in the generated annotations. Default: `default`

Example:

```bash
mammadctl convert -n production -c thr1 my-service
```

Output file:

```text
./redis-vitrin.yaml
```

### `mammadctl sotoon`

Parent command for Sotoon API utilities.

Global Sotoon options:

- `--url`: Sotoon API base URL. Default: `https://api.sotoon.ir`
- `--workspace-id`: Sotoon workspace ID. Default: `fee4bbf5-342d-4243-b4aa-f0bee508c39a`

All Sotoon commands require:

```bash
export BEPA_TOKEN="<your-bepa-token>"
```

Example with custom workspace:

```bash
mammadctl sotoon --workspace-id <workspace-id> user list
```

### `mammadctl sotoon user list`

Lists Sotoon users and service-users with their UUIDs, emails and names.

It's useful for getting the UUID of a user to use in other commands.

```bash
mammadctl sotoon user list
```

Example:

```bash
mammadctl sotoon user list | grep mohammad.jamshidi
```

### `mammadctl sotoon access cdn`

Adds `cdn-editor` and `dns-viewer` roles to the specified user for a domain.

```bash
mammadctl sotoon access cdn <user-uuid> <domain>
```

Example:

```bash
mammadctl sotoon access cdn c83786b6-dc5e-495f-89f7-04be64374cb9 example.com
```

### `mammadctl sotoon access k8s`

Adds namespaced Kubernetes access to the specified Sotoon user.

Supported datacenters:
- `neda`
- `afra`

```bash
mammadctl sotoon access k8s <user-uuid> <datacenter> <namespace>
```

Example:

```bash
mammadctl sotoon access k8s c83786b6-dc5e-495f-89f7-04be64374cb9 neda default
```

### `mammadctl sotoon access s3`

Adds full access to an S3 bucket for the specified user.

```bash
mammadctl sotoon access s3 <user-uuid> <datacenter> <bucket-name>
```

Supported datacenters:

- `neda`
- `afra`

Required environment variables:

```bash
export S3_ACCESS_KEY="<your-s3-access-key>"
export S3_SECRET_KEY="<your-s3-secret-key>"
```

Options:

- `--neda-s3-url`: Neda S3 base URL. Default: `s3.thr1.sotoon.ir`
- `--afra-s3-url`: Afra S3 base URL. Default: `s3.thr2.sotoon.ir`

Example:

```bash
mammadctl sotoon access s3 c83786b6-dc5e-495f-89f7-04be64374cb9 neda my-bucket
```
