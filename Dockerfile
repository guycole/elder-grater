# 
# elder-grater
#
# docker build . -t elder-grater:1
# kubectl run eg --rm -it --image=elder-grater:1 --restart=Never -- sh
#
FROM --platform=linux/amd64 golang:1.21.7-alpine3.19 as buildx
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