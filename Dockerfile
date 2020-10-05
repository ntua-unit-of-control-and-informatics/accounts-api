FROM golang:1.13 AS build

LABEL maintainer="Pantelis Karatzas <pantelispanka@gmail.com>"

COPY . /app
WORKDIR /app
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o accounts-api .

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /root/
# COPY accounts-api .
EXPOSE 8000
COPY --from=build /app/accounts-api .
CMD ["./accounts-api"]
