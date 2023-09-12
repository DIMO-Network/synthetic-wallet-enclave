FROM amazonlinux:2 as build

RUN yum -y install wget

RUN amazon-linux-extras install aws-nitro-enclaves-cli docker && \
    yum install aws-nitro-enclaves-cli-devel jq util-linux \
        unzip curl which file -y
        
CMD ["/bin/bash"] 


