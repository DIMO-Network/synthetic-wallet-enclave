FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli

WORKDIR /deps

COPY resources/docker/run.sh .
RUN wget "https://github.com/DIMO-Network/synthetic-wallet-enclave/releases/download/v1.7/synthetic-wallet-enclave.eif"

FROM amazonlinux:2 as image

COPY --from=build /usr/bin/nitro-cli /
COPY --from=build /deps/synthetic-wallet-enclave.eif /

RUN mkdir -p /var/log/nitro_enclaves

COPY --from=build /deps/run.sh /

CMD ["/run.sh"]
