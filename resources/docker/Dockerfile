FROM amazonlinux:2 as full_image

# Install packages
#
RUN amazon-linux-extras install aws-nitro-enclaves-cli -y && \
    yum install aws-nitro-enclaves-cli-devel curl jq -y && \
    yum clean all && \
    rm -rf /var/cache/yum && \
    mkdir -p /enclave
COPY resources/docker/run.sh  /
COPY resources/docker/synthetic-wallet-enclave.eif /eif/

CMD ["/run.sh"]
