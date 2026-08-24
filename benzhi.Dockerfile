FROM golang:1.23.12
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor ./...
ENV BMS_ADDR=:18080 BMS_DATA=/app/data BMS_SEED=true
EXPOSE 18080
CMD ["bash"]
