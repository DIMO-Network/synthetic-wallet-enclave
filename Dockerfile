FROM amazonlinux:2

WORKDIR /

COPY ./kmstool_enclave_cli ./
COPY ./libnsm.so ./
COPY ./test-enclave ./

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/"

ENTRYPOINT ["./test-enclave", "5000"]
