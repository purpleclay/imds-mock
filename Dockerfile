FROM scratch

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini-static
RUN chmod +x /tini-static

COPY imds-mock /imds-mock

ENTRYPOINT ["/tini-static", "--", "/imds-mock"]