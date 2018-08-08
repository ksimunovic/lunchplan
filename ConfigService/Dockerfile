FROM launchplan/gocompileimage AS builder
ARG foldername
LABEL maintainer ksimunovic
WORKDIR $foldername
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
ARG foldername
WORKDIR /root/
COPY --from=builder $foldername/app .
EXPOSE 50000
CMD ["./app"]