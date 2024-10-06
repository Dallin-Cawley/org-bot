# syntax=docker/dockerfile:1.2

FROM golang:1.23-alpine3.20 as build

COPY ./ /app

WORKDIR /app

RUN go build -o orgBot

# Deploy the final stage
FROM golang:1.23-alpine3.20 as final

ARG GIT_SHA="local development"
ENV GIT_SHA=${GIT_SHA}

COPY --from=build /app/orgBot  /app/orgBot
COPY --from=build /app/config /app/config

WORKDIR /app

ENTRYPOINT [ "/app/orgBot" ]