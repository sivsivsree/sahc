# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src  && go build -o ./bin/sahc -v ./cmd/main.go


# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/bin/ /app/
ENV SAHC_CONFIG service.yaml
ENTRYPOINT ./sahc