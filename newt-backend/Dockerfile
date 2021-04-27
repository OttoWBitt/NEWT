FROM golang:1.16 AS builder
COPY . /app
WORKDIR /app 
RUN go build -o Newt

FROM ubuntu:20.04
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends -qq \
        ca-certificates \
    && apt-get install -y tzdata \
    && apt-get clean --dry-run \
    && rm -rf /var/lib/apt/lists/*
RUN echo 'America/Sao_Paulo' > /etc/timezone && \
    rm /etc/localtime && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata
COPY /html /app/html
COPY --from=builder /app/Newt /app/Newt
RUN mkdir app/files
WORKDIR /app
ENTRYPOINT ["./Newt"]