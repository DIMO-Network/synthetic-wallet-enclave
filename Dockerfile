FROM amazonlinux:2

WORKDIR /

COPY ./kmstool_enclave_cli ./
COPY ./libnsm.so ./
COPY ./synthetic-wallet-enclave ./

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/"

ENTRYPOINT ["./synthetic-wallet-enclave", "5000"]
