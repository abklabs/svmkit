# SVMKit

SVMKit is an innovative Go package and Pulumi provider designed specifically for running Solana validator and RPC clients, including Anza and Firedancer. Tailored for flexibility, robustness, and ease of use, SVMKit targets various environments, including bare metal servers, virtual machines, and containerized platforms.

Our mission for SVMKit is to empower developers and institutions alike to harness the power of Solana's blockchain technology with minimal friction. By providing a robust, user-friendly platform, we strive to drive widespread adoption and foster innovation within the SVM ecosystem.

SVMKit drastically reduces the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience. By involving strategic partners and early customers, SVMKit leverages community feedback to continuously improve and adapt to real-world needs. It aims to be the go-to solution for installing and managing both permissioned and permissionless Solana clusters, with built-in mechanisms for easy bridging between them.

# SVMKit

SVMKit is an innovative, Debian-based OS distribution designed specifically for running Solana validator and RPC clients, including Anza and Firedancer. Tailored for flexibility, robustness, and ease of use, SVMKit targets various environments, including bare metal servers, virtual machines, and containerized platforms.

## Features

- **Validator Appliance:** Simplifies the operation and version control of Solana nodes. Provides community-tested configuration setups for seamless deployment and management.

- **Cluster Management:** Enables scalable, low-effort management of large validator and RPC clusters. Supports both permissioned and permissionless Solana clusters, catering to diverse user needs.

- **Developer Tooling:** Accelerates application development on the Solana Virtual Machine (SVM). Offers a suite of tools to streamline the building, launching, and maintaining of SVM forks.

- **On-Chain Permissions System:** Facilitates secure, scalable permission management for institutional customers. Ensures compliance and security in regulated environments.

- **Seamless Bridge to Mainnet:** Provides easy access to Solana Mainnet liquidity. Ensures interoperability and smooth transitions between different Solana environments.

## Objectives

- **Lowering Expertise Barriers:** SVMKit drastically reduces the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience.

- **Community-Driven Development:** By involving strategic partners and early customers, SVMKit leverages community feedback to continuously improve and adapt to real-world needs.

- **Comprehensive Management Solution:** SVMKit aims to be the go-to solution for installing and managing both permissioned and permissionless Solana clusters, with built-in mechanisms for easy bridging between them.

## Vision

Our vision for SVMKit is to empower developers and institutions alike to harness the power of Solana's blockchain technology with minimal friction. By providing a robust, user-friendly platform, we strive to drive widespread adoption and foster innovation within the SVM ecosystem.

## Installing

Follow the steps below to install SVMKit on your system:

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

6. Install SVMKit's build of the Agave validator:

```bash
sudo apt-get install zuma-agave-validator
```

=======

# Pulumi Native Provider Boilerplate

This repository is a boilerplate showing how to create and locally test a native Pulumi provider.

## Authoring a Pulumi Native Provider

This boilerplate creates a working Pulumi-owned provider named `xyz`.
It implements a random number generator that you can [build and test out for yourself](#test-against-the-example) and then replace the Random code with code specific to your provider.

### Prerequisites

Prerequisites for this repository are already satisfied by the [Pulumi Devcontainer](https://github.com/pulumi/devcontainer) if you are using Github Codespaces, or VSCode.

If you are not using VSCode, you will need to ensure the following tools are installed and present in your `$PATH`:

- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [Go 1.21](https://golang.org/dl/) or 1.latest
- [NodeJS](https://nodejs.org/en/) 14.x. We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
- [Yarn](https://yarnpkg.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [Python](https://www.python.org/downloads/) (called as `python3`). For recent versions of MacOS, the system-installed version is fine.
- [.NET](https://dotnet.microsoft.com/download)

### Build & test the boilerplate XYZ provider

1. Create a new Github CodeSpaces environment using this repository.
1. Open a terminal in the CodeSpaces environment.
1. Run `make build install` to build and install the provider.
1. Run `make gen_examples` to generate the example programs in `examples/` off of the source `examples/yaml` example program.
1. Run `make up` to run the example program in `examples/yaml`.
1. Run `make down` to tear down the example program.

### Creating a new provider repository

Pulumi offers this repository as a [GitHub template repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) for convenience. From this repository:

1. Click "Use this template".
1. Set the following options:
   - Owner: pulumi
   - Repository name: pulumi-xyz-native (replace "xyz" with the name of your provider)
   - Description: Pulumi provider for xyz
   - Repository type: Public
1. Clone the generated repository.

From the templated repository:

1. Run the following command to update files to use the name of your provider (third-party: use your GitHub organization/username):

   ```bash
   make prepare NAME=foo REPOSITORY=github.com/pulumi/pulumi-foo ORG=myorg
   ```

   This will do the following:

   - rename folders in `provider/cmd` to `pulumi-resource-{NAME}`
   - replace dependencies in `provider/go.mod` to reflect your repository name
   - find and replace all instances of the boilerplate `xyz` with the `NAME` of your provider.
   - find and replace all instances of the boilerplate `abc` with the `ORG` of your provider.
   - replace all instances of the `github.com/pulumi/pulumi-xyz` repository with the `REPOSITORY` location

#### Build the provider and install the plugin

```bash
$ make build install
```

This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### Test against the example

```bash
$ cd examples/simple
$ yarn link @pulumi/xyz
$ yarn install
$ pulumi stack init test
$ pulumi up
```

Now that you have completed all of the above steps, you have a working provider that generates a random string for you.

#### A brief repository overview

You now have:

1. A `provider/` folder containing the building and implementation logic
   1. `cmd/pulumi-resource-xyz/main.go` - holds the provider's sample implementation logic.
2. `deployment-templates` - a set of files to help you around deployment and publication
3. `sdk` - holds the generated code libraries created by `pulumi-gen-xyz/main.go`
4. `examples` a folder of Pulumi programs to try locally and/or use in CI.
5. A `Makefile` and this `README`.

#### Additional Details

This repository depends on the pulumi-go-provider library. For more details on building providers, please check
the [Pulumi Go Provider docs](https://github.com/pulumi/pulumi-go-provider).

### Build Examples

Create an example program using the resources defined in your provider, and place it in the `examples/` folder.

You can now repeat the steps for [build, install, and test](#test-against-the-example).

## Configuring CI and releases

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

## References

Other resources/examples for implementing providers:

- [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
- [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
  > > > > > > > c9c5061 (Initial commit)
