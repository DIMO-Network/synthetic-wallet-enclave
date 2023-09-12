FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli

COPY / /

WORKDIR /build

RUN nitro-cli build-enclave --docker-dir . --docker-uri synthetic-wallet-enclave --output-file synthetic-wallet-enclave.eif

FROM amazonlinux:2 as image

COPY --from=build /usr/bin/nitro-cli /
COPY --from=build /build/synthetic-wallet-enclave.eif /

COPY run.sh /

CMD ["/run.sh"]
