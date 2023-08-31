#!/bin/bash -e
# Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

readonly EIF_PATH="/synthetic-wallet-enclave.eif"
readonly ENCLAVE_CPU_COUNT=2
readonly ENCLAVE_MEMORY_SIZE=1024
readonly ENCLAVE_CID=16

main() {
    nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
        --eif-path $EIF_PATH --debug-mode --enclave-cid $ENCLAVE_CID

    sleep infinity
}

main