# 1.20.5-alpine3.18 https://hub.docker.com/layers/library/golang/alpine3.18/images/sha256-8448363817eaf10a11990b02a0109be4a83da52a9ac49061fb77080bcb9b19a8
FROM --platform=$BUILDPLATFORM golang:alpine3.18@sha256:8448363817eaf10a11990b02a0109be4a83da52a9ac49061fb77080bcb9b19a8 AS build
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd main && go build -o /out/hello .

FROM scratch AS bin
COPY --from=build /out/hello /

CMD [ "/hello" ]