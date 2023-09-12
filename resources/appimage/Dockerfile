FROM golang:1.21 as enclave-app-build

WORKDIR /src/synthetic-wallet-enclave

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /usr/bin/synthetic-wallet-enclave

FROM amazonlinux:2 AS release

WORKDIR /

COPY --from=dimozone/aws-nitro-enclaves-sdk-c:latest /usr/lib64/libnsm.so /usr/lib64/libnsm.so
COPY --from=dimozone/aws-nitro-enclaves-sdk-c:latest /usr/bin/kmstool_enclave_cli /
COPY --from=enclave-app-build /usr/bin/synthetic-wallet-enclave /

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/usr/lib64/"

USER nonroot:nonroot

EXPOSE 5000

ENTRYPOINT ["./synthetic-wallet-enclave", "5000"]
