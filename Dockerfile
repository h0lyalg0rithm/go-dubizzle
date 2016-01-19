# FROM golang:onbuild
FROM ubuntu:14.04

# gcc for cgo
RUN apt-get update && apt-get install -qq -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		pkg-config \
		make \
		curl \
		wget \
		libxml2-dev \
		libxslt1-dev \
		python-dev \
		openssh-client \
		ca-certificates \
		git \
	&& rm -rf /var/lib/apt/lists/*

# Make ssh dir
RUN mkdir /root/.ssh/

# Create known_hosts
RUN touch /root/.ssh/known_hosts

# Copy over private key, and set permissions
# ADD id_rsa /root/.ssh/id_rsa
# RUN ssh-keygen -t rsa -b 4096 -C "golang-docker"
RUN ssh-keygen -q -t rsa -N '' -f /id_rsa

# Add bitbuckets key
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN ssh-keyscan golang.org >> /root/.ssh/known_hosts
RUN ssh-keyscan storage.googleapis.com >> /root/.ssh/known_hosts

# ENV GOLANG_VERSION 1.5.3
# ENV GOLANG_DOWNLOAD_URL https://storage.googleapis.com/golang/go1.5.3.src.tar.gz
# ENV GOLANG_DOWNLOAD_SHA1 754e06dab1c31ab168fc9db9e32596734015ea9e24bc44cae7f237f417ce4efe

ENV GOLANG_VERSION 1.4.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_DOWNLOAD_SHA1 486db10dc571a55c8d795365070f66d343458c48

RUN curl -fsSLk "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA1  golang.tar.gz" | sha1sum -c - \
	&& tar -C /usr/src -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/src/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY go-wrapper /usr/local/bin/

RUN go get github.com/codegangsta/cli
RUN go get github.com/moovweb/gokogiri

# ENV GOROOT /go/src
ENV APP_HOME /go/src/go-dubizzle

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME
