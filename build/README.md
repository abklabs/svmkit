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
[remote "yellowstone-grpc"]
	url = git@github.com:rpcpool/yellowstone-grpc.git
	fetch = +refs/tags/*:refs/remotes/yellowstone-grpc/*
```

Once that's in place, you can run `solana-build` and get something like the following:

```
$ ../solana-build
INFO	git fetching remote solana-labs...
INFO	git fetching remote anza-xyz...
INFO	git fetching remote PowerLedger...
INFO	git fetching remote jito-foundation...
INFO	git fetching remote pyth-network...
INFO	building solana-labs/master -> solana-validator inside build-solana-labs/master-1448251...
Removing target/
branch 'build-solana-labs/master-1448251' set up to track 'solana-labs/master'.
Switched to a new branch 'build-solana-labs/master-1448251'
rust_nightly=nightly-2024-01-05
ci_docker_image=solanalabs/ci:rust_1.76.0_nightly-2024-01-05
rust_stable=1.76.0
.
.
.
```

Which will produce several packages in `../build` e.g.:

```
$ find ../build  -type f
../build/solana-labs/master/svmkit-solana-validator_2.0.0-1_amd64.deb
../build/jito-foundation/master/svmkit-jito-validator_2.1.0-1_amd64.deb
../build/anza-xyz/master/svmkit-agave-validator_2.1.0-1_amd64.deb
../build/pyth-network/pyth-v1.14.17/svmkit-pyth-validator_1.14.177-1_amd64.deb
../build/PowerLedger/upgrade_to_v1.16.28/svmkit-powerledger-validator_1.16.28-1_amd64.deb
