# mohammadctl

`mohammadctl` is a small CLI for Kubernetes helpers and Sotoon API access management.

## Installation

### Prerequisites

- Go `1.26.0` or newer, as defined in `go.mod`
- Access to a Kubernetes cluster for Kubernetes-related commands
- Sotoon API credentials for Sotoon-related commands

### Build from source

```bash
git clone <repository-url>
cd mohammadctl
go build -o mohammadctl .
```

You can then run it from the project directory:

```bash
./mohammadctl --help
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
mohammadctl --help
```

## Environment variables

### Required for Sotoon commands

Export this before using any `mohammadctl sotoon ...` command:

```bash
export BEPA_TOKEN="<your-bepa-token>"
```

### Required for `sotoon access s3`

Export these before using `mohammadctl sotoon access s3 ...`:

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

- `mohammadctl decode secret`
- `mohammadctl convert`

## Commands

```text
mohammadctl
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

### `mohammadctl decode secret`

Finds a Kubernetes secret and writes it to `./secret.yaml` with decoded values in `stringData` mode.

It's useful when you want to edit a `secret` or `sealedSecret`.

```bash
mohammadctl decode secret -n <namespace> <secret-name>
```

Options:

- `-n, --namespace`: Kubernetes namespace. Default: `default`

Example:

```bash
mohammadctl decode secret -n production app-secret
```

Output file:

```text
./secret.yaml
```

### `mohammadctl convert`

Reads a Kubernetes Service from the current cluster and writes a melange-compatible(Bazaar operator for mirroring services between clusters) formated Service manifest to `<namespace>-<service-name>.yaml`.

```bash
mohammadctl convert -n <namespace> -c <cluster> <service-name>
```

Options:

- `-n, --namespace`: Kubernetes namespace. Default: `default`
- `-c, --cluster`: Cluster name to put in the generated annotations. Default: `default`

Example:

```bash
mohammadctl convert -n production -c thr1 my-service
```

Output file:

```text
./redis-vitrin.yaml
```

### `mohammadctl sotoon`

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
mohammadctl sotoon --workspace-id <workspace-id> user list
```

### `mohammadctl sotoon user list`

Lists Sotoon users and service-users with their UUIDs, emails and names.

It's useful for getting the UUID of a user to use in other commands.

```bash
mohammadctl sotoon user list
```

Example:

```bash
mohammadctl sotoon user list | grep mohammad.jamshidi
```

### `mohammadctl sotoon access cdn`

Adds `cdn-editor` and `dns-viewer` roles to the specified user for a domain.

```bash
mohammadctl sotoon access cdn <user-uuid> <domain>
```

Example:

```bash
mohammadctl sotoon access cdn c83786b6-dc5e-495f-89f7-04be64374cb9 example.com
```

### `mohammadctl sotoon access k8s`

Adds namespaced Kubernetes access to the specified Sotoon user.

Supported datacenters:
- `neda`
- `afra`

```bash
mohammadctl sotoon access k8s <user-uuid> <datacenter> <namespace>
```

Example:

```bash
mohammadctl sotoon access k8s c83786b6-dc5e-495f-89f7-04be64374cb9 neda default
```

### `mohammadctl sotoon access s3`

Adds full access to an S3 bucket for the specified user.

```bash
mohammadctl sotoon access s3 <user-uuid> <datacenter> <bucket-name>
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
mohammadctl sotoon access s3 c83786b6-dc5e-495f-89f7-04be64374cb9 neda my-bucket
```
