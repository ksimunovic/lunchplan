FROM golang:latest AS builder
ARG foldername
LABEL maintainer ksimunovic
WORKDIR $foldername
COPY . .
RUN go get -u github.com/gorilla/handlers
RUN go get -u github.com/gorilla/mux
RUN go get -u golang.org/x/crypto/bcrypt
RUN go get -u gopkg.in/mgo.v2
RUN go get -u gopkg.in/mgo.v2/bson
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
ARG foldername
WORKDIR /root/
COPY --from=builder $foldername/static /root/static
COPY --from=builder $foldername/templates /root/templates
COPY --from=builder $foldername/app .
CMD ["./app"]