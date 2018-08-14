FROM jeanblanchard/alpine-glibc


RUN apk add --no-cache nodejs
RUN apk add --no-cache bash git

RUN apk add --no-cache alpine-sdk

RUN apk add --no-cache python && \
    python -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip install --upgrade pip setuptools && \
    rm -r /root/.cache

RUN apk add --no-cache zeromq-dev

RUN npm i https://github.com/qtumproject/qtumcore-node.git#master

RUN ln -s /node_modules/.bin/qtumcore-node /usr/local/bin/qtumcore-node ; \
    qtumcore-node create mynode && \
    cd mynode ; \
    qtumcore-node install https://github.com/qtumproject/insight-api.git#master


ENV QTUM_RELEASE 0.15.2
ENV QTUM_RELEASE_TAR qtum-${QTUM_RELEASE}-x86_64-linux-gnu.tar.gz

RUN wget https://github.com/qtumproject/qtum/releases/download/mainnet-ignition-v${QTUM_RELEASE}/${QTUM_RELEASE_TAR} && \
  tar -xf $QTUM_RELEASE_TAR -C /usr/local --strip-components=1 --exclude=*-qt --exclude=test_qtum --exclude=qtum-tx && \
  rm $QTUM_RELEASE_TAR


RUN cd /mynode && qtumcore-node install https://github.com/qtumproject/qtum-explorer.git#master

COPY qtumcore-node.json /etc/
RUN ln -s /usr/local/bin/qtumd /node_modules/qtumcore-node/bin/qtumd
#COPY qtum.conf /
#COPY qtumcore-node.json /root/.bitcore/qtumcore-node.json
#COPY qtum.conf /root/.bitcore/data/qtum.conf
WORKDIR /mynode
EXPOSE 3001
#CMD qtumcore-node start 
CMD qtumcore-node start -c /etc/qtumcore-node.json
