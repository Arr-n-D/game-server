FROM alpine:3.20.3 as build-valve-sockets

WORKDIR /sockets

RUN apk update && \
    apk add \
        gcc \
        g++ \
        ccache \
        cmake \
        ninja \
        pkgconf \
        git \
        linux-headers \
        go \
        protobuf-dev \
        openssl-dev 

RUN git clone --recurse-submodules -j8 https://github.com/Arr-n-D/gns.git 

RUN cd ./gns/lib/GameNetworkingSockets && \
    mkdir build && cd build && cmake -G Ninja .. && ninja

FROM alpine:3.20.3 as final

WORKDIR /app

RUN apk update && \
    apk add \
        gcc \
        g++ \
        protobuf-dev \
        openssl-dev 

COPY --from=golang:1.23-alpine /usr/local/go/ /usr/local/go/
COPY --from=build-valve-sockets /sockets/gns/lib/GameNetworkingSockets/build/bin/libGameNetworkingSockets.so /usr/lib
ENV PATH="/usr/local/go/bin:${PATH}"

COPY . /app/

RUN go build -o ./build/game-server /app/packages/server

# CMD specifies the binary to be run (with optional parameters)
CMD ["./build/game-server"]