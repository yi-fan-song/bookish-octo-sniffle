FROM golang:1.17

RUN apt-get update

RUN apt-get install -y \
    wget lsb-release

RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update

RUN apt-get -y install postgresql

WORKDIR /go/src/octo

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN chmod +x ./run.sh

RUN go build -o ./build/octo .

EXPOSE 8080

ENTRYPOINT [ "/go/src/octo/run.sh" ]