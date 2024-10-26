### Apt Setup

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
echo "deb https://apt.abklabs.com/svmkit dev main" | sudo tee /etc/apt/sources.list.d/svmkit.list
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
sudo apt-get install svmkit-agave-validator svmkit-solana-cli
```
