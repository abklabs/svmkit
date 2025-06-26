# Building Solana Validator Packages

Note: To use the scripts in this directory, you must have a recent version of the `svmkit` CLI in your path.

`solana-build` will produce a series of Debian packages out of a Solana-fork-based git repository. To build all of the forked flavors, you should configure your remotes to have all of the following in your `.git/config`:

```
[remote "solana-labs"]
	url = git@github.com:solana-labs/solana.git
	fetch = +refs/heads/*:refs/remotes/solana-labs/*
[branch "master"]
	remote = anza-xyz
	merge = refs/heads/master
[remote "anza-xyz"]
	url = git@github.com:anza-xyz/agave.git
	fetch = +refs/heads/*:refs/remotes/anza-xyz/*
[remote "pyth-network"]
	url = git@github.com:pyth-network/pythnet.git
	fetch = +refs/heads/*:refs/remotes/pyth-network/*
[remote "PowerLedger"]
	url = git@github.com:PowerLedger/powr.git
	fetch = +refs/heads/*:refs/remotes/PowerLedger/*
[remote "jito-foundation"]
	url = git@github.com:jito-foundation/jito-solana.git
	fetch = +refs/heads/*:refs/remotes/jito-foundation/*
[branch "powerledger"]
	remote = PowerLedger
	merge = refs/heads/master
[remote "mantis"]
	url = git@github.com:ComposableFi/mantis-solana.git
	fetch = +refs/heads/*:refs/remotes/mantis/*
[remote "mirrorworld-universe"]
	url = git@github.com:mirrorworld-universe/hypergrid-grid.git
	fetch = +refs/heads/*:refs/remotes/mirrorworld-universe/*
[remote "xen"]
	url = git@github.com:FairCrypto/solanalabs.git
	fetch = +refs/heads/*:refs/remotes/xen/*
[remote "tachyon"]
	url = git@github.com:x1-labs/tachyon.git
	fetch = +refs/heads/*:refs/remotes/tachyon/*
```

Once that's in place, you can run `solana-build` and get something like the following:

```
$ ../solana-build anza-xyz v2.2.17
INFO	git fetching remote anza-xyz...
INFO	building e998175 -> anza-xyz...
.
.
.
```

Which will produce several packages in `../build`.

## Building Yellowstone-GRPC

Similiarly to the above, you can build the Yellowstone-GRPC packages by adding the following to your `.git/config`:

```
[remote "yellowstone-grpc"]
	url = git@github.com:rpcpool/yellowstone-grpc.git
	fetch = +refs/tags/*:refs/remotes/yellowstone-grpc/*
```

Then you can run `yellowstone-grpc-build` to build the packages. In addition to the deps required
for the `solana-build` you will also need to install `protoc` compiler for protobufs
