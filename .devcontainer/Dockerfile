# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.19, 1.18, 1-bullseye, 1.19-bullseye, 1.18-bullseye, 1-buster, 1.19-buster, 1.18-buster
ARG VARIANT=1-bullseye
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 18, 16, 14
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
ARG SERF_VERSION="0.8.2"
RUN sudo apt-get update \ 
    && sudo apt-get install -y --no-install-recommends protobuf-compiler libprotobuf-dev \
    && sudo rm -rf /var/lib/apt/lists/* \
    && sudo curl -L -o serf.zip https://releases.hashicorp.com/serf/${SERF_VERSION}/serf_${SERF_VERSION}_linux_amd64.zip \
    && sudo unzip serf.zip -d /usr/local/serf \
    && sudo rm -f serf.zip
ENV PATH="$PATH:/usr/local/serf"

# go getは、Googleのサーバーへの接続がブロックされたため
# RUN go install google.golang.org/grpc@v1.53.0
RUN git clone https://github.com/grpc/grpc-go.git /go/src/google.golang.org/grpc

RUN go install github.com/sqs/goreturns@latest

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
RUN go install google.golang.org/protobuf/...@v1.28.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
