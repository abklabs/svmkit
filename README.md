# SVMKit

SVMKit is a comprehensive toolkit designed for deploying and managing various modules of the Solana Virtual Machine (SVM), such as consensus, genesis, and RPC. It offers operators a consistent administrative experience and provides SVM developers with clear methods for integrating their flavor of SVM module into SVMKit.

## Mission

Our mission is to empower developers and institutions to harness the power of Solana's blockchain technology with minimal friction. By offering a robust, user-friendly platform, we aim to drive widespread adoption and foster innovation within the SVM ecosystem.

## Value

SVMKit offers significant value to both software operators and SVM developers:

- **For Software Operators:** SVMKit simplifies the management of all SVM software components and variants using a single operations toolkit. This reduces the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience.

- **For SVM Developers:** SVMKit allows developers to quickly and accurately represent their software, enabling operators to leverage it using a familiar toolkit. By involving strategic partners and early customers, SVMKit continuously improves and adapts to real-world needs through community feedback.

## Goals

- **Simplify Software Operations:** Provide simplified creation and management for all modules of the Solana Virtual Machine (SVM) with community-tested configuration setups for seamless deployment and management.

- **Scalable Cluster Management:** Enable scalable, low-effort management of large validator and RPC clusters, supporting both permissioned and permissionless Solana clusters to cater to diverse user needs.

- **Easy Onboarding:** Provide a framework for representing SVM modules, enabling SVM developers the ability to integrate their solutions into SVMKit. This streamlines the process of building, launching, and maintaining SVM operators by operators.

- **Lower Expertise Barriers:** Drastically reduce the technical expertise required to deploy and manage Solana nodes, making blockchain technology more accessible to a broader audience.

- **Community-Driven Improvement:** Leverage community feedback by involving strategic partners and early customers to continuously improve and adapt SVMKit to real-world needs.

- **Comprehensive Management Solution:** SVMKit manages all components of SVM by structuring the requirements of different modular SVM components and providing clear playbooks for representing those components.

## Offerings

### Apt Repository

ABK Labs provides a build and release service for SVM software at `apt.abklabs.com`, which allows operators to install various SVM modules through apt. More details can be found in the [build](/build) directory.

Follow these steps to install the SVMKit apt repository on your system:

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

## GoLang Library

The SVMKit GoLang library, available in the [`pkg`](/pkg) directory, delivers essential tools for managing SVM modules. It includes utilities for setting up Solana genesis ledgers, configuring validators, and handling other components through SSH connections.

## Pulumi Provider

The SVMKit Pulumi provider, located in the [`provider`](/provider/) directory, facilitates infrastructure as code (IaC) for managing Solana nodes and related resources. With this provider, users can easily define, deploy, and manage their Solana infrastructure using Pulumi.

## Development

This section provides all the necessary information to start contributing to or building SVMKit. It is divided into several subsections for clarity.

#### Requirement

Ensure the following tools are installed and available in your `$PATH`:

- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [Go 1.22](https://golang.org/dl/) or 1.latest
- [NodeJS](https://nodejs.org/en/) 14.x (We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations)
- [Yarn](https://yarnpkg.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [Python](https://www.python.org/downloads/) (referred to as `python3`; the system-installed version is sufficient for recent MacOS versions)
- [.NET](https://dotnet.microsoft.com/download)

#### Build

```bash
$ make build install
```

This will build the pulumi provider, generate language sdks, and prepare host to execute the plugin locally.

#### Demo

You can find a catalog of example Pulumi projects to help you get started with SVMkit [here](./examples).

```bash
$ cd examples/aws-agave-validator
$ yarn link @pulumi/svmkit
$ yarn install
$ pulumi stack init demo
$ pulumi up
```

In this example, an Agave validator is installed on a machine via SSH, joining the Solana testnet.

Teams can add more validator clients to SVMkit, which will be accessible through the `validator` namespace in `@pulumi/svmkit`.

```typescript
new svmkit.validator.Agave(
  "validator",
  {
    connection,
    keyPairs: {
      identity: validatorKey.json,
      voteAccount: voteAccountKey.json,
    },
    flags: {
      entryPoint: [
        "entrypoint.testnet.solana.com:8001",
        "entrypoint2.testnet.solana.com:8001",
        "entrypoint3.testnet.solana.com:8001",
      ],
      knownValidator: [
        "5D1fNXzvv5NjV1ysLjirC4WY92RNsVH18vjmcszZd8on",
        "7XSY3MrYnK8vq693Rju17bbPkCN3Z7KvvfvJx4kdrsSY",
        "Ft5fbkqNa76vnsjYNwjDZUXoTWpP7VYm3mtsaQckQADN",
        "9QxCLckBiJc783jnMvXZubK4wH86Eqqvashtrwvcsgkv",
      ],
      expectedGenesisHash: "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY",
      useSnapshotArchivesAtStartup: "when-newest",
      rpcPort: 8899,
      privateRPC: true,
      onlyKnownRPC: true,
      dynamicPortRange: "8002-8020",
      gossipPort: 8001,
      rpcBindAddress: "0.0.0.0",
      walRecoveryMode: "skip_any_corrupted_record",
      limitLedgerSize: 50000000,
      blockProductionMethod: "central-scheduler",
      fullSnapshotIntervalSlots: 1000,
      noWaitForVoteToStartLeader: true,
    },
  },
  {
    dependsOn: [instance],
  }
);
```

#### Structure

The repository includes the following:

| Directory/File          | Description                                                         |
| ----------------------- | ------------------------------------------------------------------- |
| `provider/`             | Build and implementation logic for the SVMkit Pulumi provider.      |
| `pkg`                   | SVMKit Go packages.                                                 |
| `sdk`                   | Generated code libraries created by `pulumi-gen-svm/main.go`.       |
| `examples`              | Pulumi programs for local testing and CI usage.                     |
| `build`                 | Scripts for building and publishing to the AKB Labs apt repository. |
| `Makefile` and `README` | Standard project files for building and documentation.              |
