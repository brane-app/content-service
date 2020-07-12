FROM golang:latest

WORKDIR /build
COPY . .
RUN go get -u ./...
RUN go build -o built .

ENV MONKEBASE_CONNECTION ""
ENV FERROTHORN_SECRET ""
ENV FERROTHORN_HOST ""
ENTRYPOINT ./built
