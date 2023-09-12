#!/bin/bash -e
# Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

readonly EIF_PATH="/eif/synthetic-wallet-enclave.eif"
readonly ENCLAVE_CPU_COUNT=2
readonly ENCLAVE_MEMORY_SIZE=1024
readonly ENCLAVE_CID=16
readonly AWS_REGION=us-east-2
main() {
    ls -al /
    ls -al /var/log

    /nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
        --eif-path $EIF_PATH --debug-mode --enclave-cid $ENCLAVE_CID

    vsock-proxy 8000 kms.$AWS_REGION.amazonaws.com 443 &
    sleep infinity
}

main