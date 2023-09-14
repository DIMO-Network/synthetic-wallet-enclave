#!/bin/bash -e

set -x	
readonly ENCLAVE_NAME="synthetic-wallet-enclave"
readonly EIF_PATH="/eif/synthetic-wallet-enclave.eif"


ENCLAVE_CPU_COUNT=${ENCLAVE_CPU_COUNT:-1}
ENCLAVE_MEMORY_SIZE=${ENCLAVE_MEMORY_SIZE:-768}
ENCLAVE_CID=${ENCLAVE_CID:-16}
AWS_REGION=${AWS_REGION:-us-east-1}

term_handler() {
  echo 'Shutting down enclave'
  nitro-cli terminate-enclave --enclave-name $ENCLAVE_NAME
  kill -SIGTERM $(pgrep nitro-cli)
  kill -SIGTERM $(pgrep vsock-proxy)
  echo 'Shutdown complete'
  exit 0;
}

# on callback, kill the last background process, which is `tail -f /dev/null` and execute the specified handler
trap 'kill ${!}; term_handler' SIGTERM

# run application

vsock-proxy 8000 kms.$AWS_REGION.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
    --eif-path $EIF_PATH --enclave-cid $ENCLAVE_CID --attach-console &

if [[ ! -v "${ENCLAVE_DEBUG_MODE}" ]]; then
  echo 'Starting production enclave.'
  nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
    --eif-path $EIF_PATH --enclave-cid $ENCLAVE_CID &
  echo 'Enclave running.'
else
  echo 'Starting development enclave.'
  nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
    --eif-path $EIF_PATH --enclave-cid $ENCLAVE_CID --attach-console &
  echo 'Enclave started in debug mode.'
fi

# wait forever
while true
do
  tail -f /dev/null & wait ${!}
done
