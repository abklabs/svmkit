import * as pulumi from "@pulumi/pulumi";
import * as svmkit from "@pulumi/svmkit";
import * as aws from "@pulumi/aws";
import * as tls from "@pulumi/tls";

const config = new pulumi.Config("svmkit");

const ami = config.require("ami");

const sshKey = new tls.PrivateKey("ssh-key", {
  algorithm: "ED25519",
});

const keyPair = new aws.ec2.KeyPair("keypair", {
  publicKey: sshKey.publicKeyOpenssh,
});

const validatorKey = new svmkit.KeyPair("validator-key");
const voteAccountKey = new svmkit.KeyPair("vote-account-key");

const securityGroup = new aws.ec2.SecurityGroup("security-group", {
  description: "Allow SSH and specific inbound traffic",
  ingress: [
    {
      protocol: "tcp",
      fromPort: 22,
      toPort: 22,
      cidrBlocks: ["0.0.0.0/0"],
    },
    {
      protocol: "tcp",
      fromPort: 8000,
      toPort: 8020,
      cidrBlocks: ["0.0.0.0/0"],
    },
    {
      protocol: "udp",
      fromPort: 8000,
      toPort: 8020,
      cidrBlocks: ["0.0.0.0/0"],
    },
  ],
  egress: [
    {
      protocol: "-1",
      fromPort: 0,
      toPort: 0,
      cidrBlocks: ["0.0.0.0/0"],
    },
  ],
});

const instance = new aws.ec2.Instance("instance", {
  ami,
  instanceType: "m5.2xlarge",
  keyName: keyPair.keyName,
  vpcSecurityGroupIds: [securityGroup.id],
  rootBlockDevice: {
    volumeSize: 500,
    volumeType: "gp3",
    iops: 3000,
    throughput: 250,
  },
  ebsBlockDevices: [
    {
      deviceName: "/dev/sdf",
      volumeSize: 500,
      volumeType: "gp3",
      iops: 5000,
      throughput: 700,
    },
    {
      deviceName: "/dev/sdg",
      volumeSize: 2048,
      volumeType: "gp3",
      iops: 10000,
      throughput: 700,
    },
  ],
  userData: `#!/bin/bash
    mkfs -t ext4 /dev/sdf
    mkfs -t ext4 /dev/sdg
    mkdir -p /home/sol/accounts
    mkdir -p /home/sol/ledger
    mount /dev/sdf /home/sol/accounts
    mount /dev/sdg /home/sol/ledger
  `,
});

const connection = {
  host: instance.publicDns,
  user: "admin",
  privateKey: sshKey.privateKeyOpenssh,
};

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

export const PUBLIC_DNS_NAME = instance.publicDns;
export const SSH_PRIVATE_KEY = sshKey.privateKeyOpenssh;
