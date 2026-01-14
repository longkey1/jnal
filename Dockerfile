FROM gcr.io/distroless/static:nonroot

ARG TARGETARCH
COPY jnal_linux_${TARGETARCH} /usr/local/bin/jnal

ENV JNAL_CONFIG=/app/config.toml
WORKDIR /app

ENTRYPOINT ["jnal"]
