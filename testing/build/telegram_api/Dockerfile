FROM alpine:3.14

RUN  apk update \
    && apk upgrade \
    && apk add alpine-sdk linux-headers git zlib-dev openssl-dev gperf cmake

WORKDIR /api

RUN git clone --recursive https://github.com/tdlib/telegram-bot-api.git

WORKDIR /api/telegram-bot-api

RUN rm -rf build \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=.. .. \
    && cmake --build . --target install

ENTRYPOINT ["tail"]
CMD ["-f", "/dev/null"]
