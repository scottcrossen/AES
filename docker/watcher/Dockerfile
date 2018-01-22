FROM golang

RUN apt-get update \
  && apt-get install -y inotify-tools curl golang-glide \
  && rm -rf /var/lib/apt/lists/*

COPY docker-entrypoint.sh /docker-entrypoint.sh

WORKDIR /go/src

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD []
