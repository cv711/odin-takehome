FROM golang:1.24.3-bullseye AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN --mount=type=cache,target="/root/.cache/go-build" \
  go mod download
COPY server ./server
ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
  go build -o a.out ./server
RUN chmod +x ./server

FROM scratch
EXPOSE 8080
WORKDIR /
ENV IS_PROD=true
COPY --from=build /build/a.out /server
ENTRYPOINT [ "/server" ]
