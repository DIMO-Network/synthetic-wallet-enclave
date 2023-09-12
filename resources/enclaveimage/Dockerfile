FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli docker && \
    yum install wget git aws-nitro-enclaves-cli-devel -y


WORKDIR /build/
RUN nitro-cli describe-enclaves 
RUN docker pull docker.io/dimozone/synthetic-wallet-enclave-app:local
RUN nitro-cli build-enclave --docker-uri docker.io/dimozone/synthetic-wallet-enclave-app:local --output-file synthetic-wallet-enclave.eif