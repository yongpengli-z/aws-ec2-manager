FROM golang:1.17 as builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine
WORKDIR /app
RUN apk update && apk add tzdata
RUN mkdir -p /app/docs
COPY --from=builder  /app/aws-ec2-manager /app
COPY --from=builder  /app/docs/* /app/docs
