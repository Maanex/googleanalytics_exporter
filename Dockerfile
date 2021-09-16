FROM golang:alpine
LABEL description="Obtains Google Analytics RealTime API metrics, and presents them to prometheus for scraping."

RUN mkdir /gaexporter
WORKDIR /gaexporter

COPY . .

#Install Glide, Git and dependencies
RUN apk --update add git openssh && \
    apk add --update ca-certificates && \
    apk add --no-cache curl && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

RUN go build ganalytics.go

CMD go run ganalytics.go
