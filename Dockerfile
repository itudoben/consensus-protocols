# 1.20.0-alpine3.17
FROM --platform=$BUILDPLATFORM golang:alpine3.17@sha256:405962195c7fd525604cb74ab86cb7c88fcfc30af0e31a5b3c0636a7c4e9e567 AS build
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} cd main && go build -o /out/hello .

FROM scratch AS bin
COPY --from=build /out/hello /

CMD [ "/hello" ]