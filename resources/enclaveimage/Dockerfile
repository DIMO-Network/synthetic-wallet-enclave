FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli


WORKDIR /build/

RUN nitro-cli build-enclave --docker-uri synthetic-wallet-enclave-app:latest --output-file synthetic-wallet-enclave.eif