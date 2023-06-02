On the build box:

```shell
nitro-cli build-enclave --docker-dir . --docker-uri synthetic-wallet-enclave --output-file syn
thetic-wallet-enclave.eif
```

On the deployment:

```shell
sudo nitro-cli run-enclave --eif-path synthetic-wallet-enclave.eif  --cpu-count 2 --memory 1024 --debug-mode --enclave-cid 16
```
