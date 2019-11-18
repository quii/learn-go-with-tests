FROM golang:1.13.4
MAINTAINER EngHyu <roomedia@naver.com>

RUN git clone https://github.com/Gunforge/learn-go-with-tests-ko /root/workspace
RUN go get golang.org/x/tools/cmd/godoc
RUN echo 'export GOPATH=/root/workspace:/go' >> ~/.bashrc
RUN mkdir /root/workspace/src
RUN mkdir /root/workspace/src/your_code
RUN mkdir /root/workspace/src/your_code/write_down_here
RUN cp /root/workspace/hello-world/v1/hello.go /root/workspace/src/your_code/write_down_here

WORKDIR /root/workspace
