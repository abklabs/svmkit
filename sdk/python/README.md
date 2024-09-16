# SVMKit

SVMKit is a comprehensive toolchain for deploying and managing Solana Virtual Machine (SVM) validators and networks. It features a Go package with interfaces that enable developers to integrate support for their specific validator clients or networks. This allows for consistent and managed operation of SVM software by validator operators.

## Mission

Our mission for SVMKit is to empower developers and institutions alike to harness the power of Solana's blockchain technology with minimal friction. By providing a robust, user-friendly platform, we strive to drive widespread adoption and foster innovation within the SVM ecosystem.

## Vision

Our vision for SVMKit is to empower developers and institutions alike to harness the power of Solana's blockchain technology with minimal friction. By providing a robust, user-friendly platform, we strive to drive widespread adoption and foster innovation within the SVM ecosystem.

## Value

SVMKit drastically reduces the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience. By involving strategic partners and early customers, SVMKit leverages community feedback to continuously improve and adapt to real-world needs. It aims to be the go-to solution for installing and managing both permissioned and permissionless Solana clusters.

## Features

- **Validator Appliance:** Simplifies the operation and version control of Solana nodes. Provides community-tested configuration setups for seamless deployment and management.

- **Cluster Management:** Enables scalable, low-effort management of large validator and RPC clusters. Supports both permissioned and permissionless Solana clusters, catering to diverse user needs.

- **Developer Tooling:** Accelerates application development on the Solana Virtual Machine (SVM). Offers a suite of tools to streamline the building, launching, and maintaining of SVM forks.

## Objectives

- **Lowering Expertise Barriers:** SVMKit drastically reduces the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience.

- **Community-Driven Development:** By involving strategic partners and early customers, SVMKit leverages community feedback to continuously improve and adapt to real-world needs.

- **Comprehensive Management Solution:** SVMKit aims to be the go-to solution for installing and managing both permissioned and permissionless Solana clusters, with built-in mechanisms for easy bridging between them.

## Apt Repository

Follow the steps below to install SVMKit apt repository on your system:

1. Update your package lists:

```bash
sudo apt-get update
```

2. Install the necessary prerequisites:

```bash
sudo apt-get install -y curl gnupg
```

3. Add the SVMKit repository to your system's software repository list:

```bash
echo "deb https://apt.abklabs.com/zuma dev main" | sudo tee /etc/apt/sources.list.d/zuma.list
```

4. Import the repository's GPG key:

```bash
curl -s https://apt.abklabs.com/keys/abklabs-archive-dev.asc | sudo apt-key add -
```

5. Update your package lists again:

```bash
sudo apt-get update
```

6. Install SVMKit's build of the Agave validator and solana cli:

```bash
sudo apt-get install zuma-agave-validator zuma-solana-cli
```

### Dependencies

You will need to ensure the following tools are installed and present in your `$PATH`:

- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [Go 1.21](https://golang.org/dl/) or 1.latest
- [NodeJS](https://nodejs.org/en/) 14.x. We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
- [Yarn](https://yarnpkg.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [Python](https://www.python.org/downloads/) (called as `python3`). For recent versions of MacOS, the system-installed version is fine.
- [.NET](https://dotnet.microsoft.com/download)

#### Build the Pulumi provider and install the plugin

```bash
$ make build install
```

#### Demo

```bash
$ cd examples/aws-agave-validator
$ yarn link @pulumi/svmkit
$ yarn install
$ pulumi stack init demo
$ pulumi up
```

#### Repository Overview

The repository includes the following:

- `provider/`: Contains the build and implementation logic for the Pulumi provider.
- `cmd/pulumi-resource-svmkit/main.go`: Contains the sample implementation logic for the provider.
- `pkg`: Contains the SVMKit Go packages.
- `sdk`: Contains the generated code libraries created by `pulumi-gen-svm/main.go`.
- `examples`: Contains Pulumi programs for local testing and CI usage.
- `build`: Contains scripts for building and publishing to the AKB Labs apt repository.
- `Makefile` and `README`: Standard project files for building and documentation.
