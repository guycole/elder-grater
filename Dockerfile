# 
# elder-grater
#
# docker build . -t elder-grater:1
# kubectl run eg --rm -it --image=elder-grater:1 --restart=Never -- sh
#
#FROM golang:1.21.4-alpine3.17
FROM --platform=linux/amd64 golang:1.21.4-alpine3.17 as buildx
#
ENV AWS_WEB_IDENTITY_TOKEN_FILE "/var/run/secrets/eks.amazonaws.com/serviceaccount/token"
#
WORKDIR /app
#
COPY go.mod .
COPY go.sum .
RUN go mod download
#
COPY *.go ./
#
RUN go build -o /app/elder-grater
#
ENTRYPOINT [ "/app/elder-grater" ]
CMD ["bogus"]
#