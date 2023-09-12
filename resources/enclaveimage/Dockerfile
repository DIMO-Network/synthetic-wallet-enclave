FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli && \
    yum install wget git aws-nitro-enclaves-cli-devel -y


WORKDIR /build/
RUN nitro-cli describe-enclaves 
RUN nitro-cli build-enclave --docker-uri synthetic-wallet-enclave-app:test --output-file synthetic-wallet-enclave.eif