# 1.21.3-alpine3.18 https://hub.docker.com/layers/library/golang/1.21.3-alpine3.18/images/sha256-533470173383661c84e0b4d9b1f806b7caeae4b3cee918cdf49fb558eee195c4?context=explore
FROM --platform=$BUILDPLATFORM golang:1.21.3-alpine3.18@sha256:533470173383661c84e0b4d9b1f806b7caeae4b3cee918cdf49fb558eee195c4 AS build
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd main && go build -o /out/hello .

FROM scratch AS bin
COPY --from=build /out/hello /

CMD [ "/hello" ]