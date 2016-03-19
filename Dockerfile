FROM ubuntu:15.10

MAINTAINER n3wtron@gmail.com

RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git mercurial bzr libgit2-dev libldap2-dev pkg-config
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1
RUN mkdir /gopath

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN go get -v -d github.com/gitDashboard/gitDashboard | echo 0
RUN go get -v github.com/revel/cmd/revel
RUN go get -v ... |echo 0
RUN go get -v ... |echo 0
RUN /gopath/bin/revel build github.com/gitDashboard/gitDashboard /gitDashboard

EXPOSE 9000

ENTRYPOINT ["/gitDashboard/run.sh"]

