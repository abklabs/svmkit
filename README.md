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

For detailed setup instructions to install the SVMKit apt repository on your system, please refer to the [APT-SETUP.md](APT-SETUP.md) file.

```bash
sudo apt-get install svmkit-agave-validator svmkit-solana-cli
```

## GoLang Library

The SVMKit GoLang library, available in the [`pkg`](/pkg) directory, delivers essential tools for managing SVM modules. It includes utilities for setting up Solana genesis ledgers, configuring validators, and handling other components through SSH connections.

```bash
go get github.com/abklabs/svmkit/pkg
```

## Pulumi Provider

The SVMKit Pulumi provider facilitates infrastructure as code (IaC) for managing Solana nodes and related resources. For more details and usage instructions, please see the [Pulumi SVMKit repository](https://github.com/abklabs/pulumi-svmkit).

```typescript
import * as svmkit from "@pulumi/svmkit";

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

## Development

This section provides all the necessary information to start contributing to or building SVMKit. It is divided into several subsections for clarity.

#### Requirement

Ensure the following tools are installed and available in your `$PATH`:

- [Go 1.22](https://golang.org/dl/) or 1.latest
- [`golangci-lint`](https://golangci-lint.run/install)

#### Structure

The repository includes the following:

| Directory/File | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `pkg`          | SVMKit Go packages.                                           |
| `build`        | Scripts for building validators and distributing through APT. |
| `README`       | Standard project files for building and documentation.        |
| `Makefile`     | Contains commands for testing and building Go packages.       |
