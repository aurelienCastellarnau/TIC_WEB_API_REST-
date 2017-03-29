FROM debian
WORKDIR /tmp
ADD . /go/src/rest/api
COPY bdd.sql /tmp
RUN export TERM=xterm && \
    export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y golang git mysql-server mysql-client
ENV GOPATH /go
RUN go get github.com/gorilla/mux github.com/go-sql-driver/mysql github.com/dgrijalva/jwt-go
CMD /etc/init.d/mysql start && \
    mysql -u root -e "create database rest" && \
    mysql -u root rest < bdd.sql && \
    go run /go/src/rest/api/*.go
EXPOSE 3000