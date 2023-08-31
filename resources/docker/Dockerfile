FROM amazonlinux:2 as build

WORKDIR /deps

COPY resources/docker/run.sh .
RUN wget "https://github.com/DIMO-Network/synthetic-wallet-enclave/releases/download/v1.7/synthetic-wallet-enclave.eif"

FROM amazonlinux:2 as image

COPY --from=build synthetic-wallet-enclave.eif /
COPY --from=build run.sh /
