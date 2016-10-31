FROM golang:1.7.3

ENV GLIDE_VERSION=v0.12.3

# install glide
RUN curl -sSL https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.tar.gz | tar -vxz -C /usr/local/bin --strip=1

COPY . /go/src/github.com/fishworks/api
WORKDIR /go/src/github.com/fishworks/api

RUN glide install

RUN make build && mv bin/api /bin

EXPOSE 8080
CMD /bin/api
