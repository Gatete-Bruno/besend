FROM golang:1.21 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY api/ api/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o manager main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
EXPOSE 8080 8081
ENTRYPOINT ["/manager"]
