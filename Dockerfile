# 1.20.0-alpine3.17
FROM --platform=$BUILDPLATFORM golang@sha256:ebceb16dc094769b6e2a393d51e0417c19084ba20eb8967fb3f7675c32b45774 AS build
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/hello .

FROM scratch AS bin
COPY --from=build /out/hello /

CMD [ "/hello" ]