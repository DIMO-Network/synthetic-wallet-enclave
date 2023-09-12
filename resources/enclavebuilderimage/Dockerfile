FROM amazonlinux:2 as build

RUN yum -y install wget
RUN amazon-linux-extras install aws-nitro-enclaves-cli docker && \
    yum install wget git aws-nitro-enclaves-cli-devel  -y

CMD ["/bin/bash"] 


