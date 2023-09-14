FROM golang:1.21 as enclave-app-build


ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 enclave

WORKDIR /src/synthetic-wallet-enclave

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /usr/bin/synthetic-wallet-enclave

FROM debian:bullseye-slim AS release

WORKDIR /

COPY --from=dimozone/aws-nitro-enclaves-sdk-c:latest /usr/lib64/libnsm.so /usr/lib64/libnsm.so
COPY --from=dimozone/aws-nitro-enclaves-sdk-c:latest /usr/bin/kmstool_enclave_cli /
COPY --from=enclave-app-build /usr/bin/synthetic-wallet-enclave /

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/usr/lib64/"

USER enclave

EXPOSE 5000

ENTRYPOINT ["./synthetic-wallet-enclave", "5000"]
