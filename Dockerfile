FROM busybox:1.35.0 AS build

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static-amd64 /tini-static-amd64
RUN chmod +rwx /tini-static-amd64

FROM scratch

COPY --from=build /tini-static-amd64 /bin/tini-static-amd64
COPY imds-mock /bin/imds-mock

ENTRYPOINT ["tini-static-amd64", "--", "imds-mock"]