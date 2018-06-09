FROM golang:latest AS builder
ARG foldername
LABEL maintainer ksimunovic
WORKDIR $foldername/ApiGateway
COPY . .
RUN go get -u github.com/gorilla/handlers
RUN go get -u github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
ARG foldername
WORKDIR /root/
COPY --from=builder $foldername/ApiGateway/certs /root/certs
COPY --from=builder $foldername/ApiGateway/app .
EXPOSE 4430
EXPOSE 80
CMD ["./app"]