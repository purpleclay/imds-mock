FROM busybox:1.36.0 AS build

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini-static
RUN chmod +rwx /tini-static

FROM scratch

# Explicitly turn on release mode within Gin
ENV GIN_MODE=release

COPY --from=build /tini-static /bin/tini-static
COPY imds-mock /bin/imds-mock

ENTRYPOINT ["tini-static", "--", "imds-mock"]