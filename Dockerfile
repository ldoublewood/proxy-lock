FROM golang:1.13.4

ENV VERSION 1.0

WORKDIR /workdir

RUN mkdir src
COPY go.mod src/
COPY go.sum src/

RUN cd src && go mod download

COPY . src

RUN cd src && go build -o ../proxylock main.go && cd ..

RUN rm -rf src

EXPOSE 3344

CMD ["/workdir/proxylock"]
