FROM qtum/qtum-bitcore

RUN apt-get update
#RUN apt-get update \
#    && \
RUN    apt-get install -y -qq --no-install-recommends \
      curl \
      bash \
      python

#    && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://deb.nodesource.com/setup_6.x | bash
RUN apt-get install -y nodejs

RUN npm i https://github.com/qtumproject/qtumcore-node.git#master

RUN ln -s /node_modules/.bin/qtumcore-node /usr/local/bin/qtumcore-node ; \
    qtumcore-node create mynode

RUN cd mynode \
    && qtumcore-node install https://github.com/qtumproject/insight-api.git#master \
    && qtumcore-node install https://github.com/qtumproject/qtum-explorer.git#master

RUN wget -v https://github.com/qtumproject/solar/releases/download/0.0.14/solar-linux-amd64 -O /usr/local/bin/solar && chmod 0755 /usr/local/bin/solar 

COPY qtumcore-node.json /etc/

RUN ln -s /usr/local/bin/qtumd /node_modules/qtumcore-node/bin/qtumd

#COPY qtum.conf /
#COPY qtumcore-node.json /root/.bitcore/qtumcore-node.json
#COPY qtum.conf /root/.bitcore/data/qtum.conf
WORKDIR /mynode
EXPOSE 3001
EXPOSE 18332
#CMD qtumcore-node start 
CMD qtumcore-node start -c /etc/qtumcore-node.json
