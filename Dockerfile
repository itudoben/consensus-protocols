# 1.20.6-alpine3.18 https://hub.docker.com/layers/library/golang/alpine3.18/images/sha256-8bdf832c26fff72bca1fa6683b1f01dd6462b2bd665aab683f9bd14bc2098b38?context=explore
FROM --platform=$BUILDPLATFORM golang:alpine3.18@sha256:8bdf832c26fff72bca1fa6683b1f01dd6462b2bd665aab683f9bd14bc2098b38 AS build
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd main && go build -o /out/hello .

FROM scratch AS bin
COPY --from=build /out/hello /

CMD [ "/hello" ]