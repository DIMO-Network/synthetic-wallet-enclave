#!/bin/bash -e
# Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

set -x	
readonly EIF_PATH="/eif/synthetic-wallet-enclave.eif"

if [[ -z "${CPU_LIMIT}" ]]; then
  ENCLAVE_CPU_COUNT=2
else
  ENCLAVE_CPU_COUNT="${CPU_LIMIT}"
fi

if [[ -z "${HUGEPAGES_LIMIT}" ]]; then
  ENCLAVE_MEMORY_SIZE=512
else
  ENCLAVE_MEMORY_SIZE="${HUGEPAGES_LIMIT}"
fi

if [[ -z "${AWS_REGION}" ]]; then
  AWS_REGION=us-east-1
else
  AWS_REGION="${AWS_REGION}"
fi

if [[ -z "${ENCLAVE_CID}" ]]; then
  ENCLAVE_CID=16
else
  ENCLAVE_CID="${ENCLAVE_CID}"
fi

main() {
    nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
        --eif-path $EIF_PATH --debug-mode --enclave-cid $ENCLAVE_CID

    vsock-proxy 8000 kms.$AWS_REGION.amazonaws.com 443 &
    sleep infinity
}

main

#nitro-cli terminate-enclave --all
#nitro-cli describe-enclaves | jq -r .[0].EnclaveID