FROM golang
ENV GOBIN=/go/bin
ADD . /go/src/github.com/Snorlock/shoppingApi
RUN go get -u github.com/dancannon/gorethink && \
    go get -u github.com/dgrijalva/jwt-go && \
    go get -u github.com/markbates/goth && \
    go get -u github.com/Snorlock/mux && \
    go get -u github.com/gorilla/sessions && \
    go get -u golang.org/x/oauth2 && \
    go get -u github.com/pborman/uuid
RUN go install src/github.com/Snorlock/shoppingApi/*.go
ENTRYPOINT /go/bin/main

EXPOSE 8000
