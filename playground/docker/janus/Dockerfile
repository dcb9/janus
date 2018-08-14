FROM alpine:latest  

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl

WORKDIR /root/

ENV QTUM_RPC_USER qtum
ENV QTUM_RPC_PASS test

EXPOSE 23889

COPY janus-linux-amd64 /usr/local/bin/janus

CMD janus --qtum-rpc="http://${QTUM_RPC_USER}:${QTUM_RPC_PASS}@qtumd:18332" --bind=0.0.0.0 --port=23889 --dev
