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

// Faucet for receiving SOL
const faucetKey = new svmkit.KeyPair("faucet-key");

// Treasury used to distribute stake
const treasuryKey = new svmkit.KeyPair("treasury-key");

// Bootstrap node
const identityKey = new svmkit.KeyPair("identity-key");
const voteAccountKey = new svmkit.KeyPair("vote-account-key");
const stakeAccountKey = new svmkit.KeyPair("stake-account-key");

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
    {
      protocol: "tcp",
      fromPort: 8899,
      toPort: 8899,
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
  instanceType: "m5.large",
  keyName: keyPair.keyName,
  vpcSecurityGroupIds: [securityGroup.id],
  rootBlockDevice: {
    volumeSize: 250,
  },
});

const connection = {
  host: instance.publicDns,
  user: "admin",
  privateKey: sshKey.privateKeyOpenssh,
};

const genesis = pulumi
  .all([
    identityKey.publicKey,
    voteAccountKey.publicKey,
    stakeAccountKey.publicKey,
    faucetKey.publicKey,
    treasuryKey.publicKey,
  ])
  .apply(
    ([
      identityPubkey,
      votePubkey,
      stakePubkey,
      faucetPubkey,
      treasuryPubkey,
    ]) => {
      const primordial = [
        {
          pubkey: identityPubkey,
          lamports: "10000000000", // 100 SOL
        },
        {
          pubkey: treasuryPubkey,
          lamports: "100000000000000", // 100000 SOL
        },
        {
          pubkey: faucetPubkey,
          lamports: "1000000000000", // 1000 SOL
        },
      ];

      return new svmkit.genesis.Solana(
        "genesis",
        {
          connection,
          flags: {
            ledgerPath: "/home/sol/ledger",
            identityPubkey,
            votePubkey,
            stakePubkey,
            faucetPubkey,
          },
          primordial,
        },
        { dependsOn: [instance] }
      );
    }
  );

new svmkit.validator.Agave(
  "validator",
  {
    connection,
    keyPairs: {
      identity: identityKey.json,
      voteAccount: voteAccountKey.json,
    },
    flags: {
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
    dependsOn: [genesis],
  }
);

export const GENESIS_HASH = genesis.genesisHash;
export const PUBLIC_DNS_NAME = instance.publicDns;
export const SSH_PRIVATE_KEY = sshKey.privateKeyOpenssh;
