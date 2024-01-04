ARG TARGETOS
ARG TARGETARCH

FROM alpine:3.19

RUN mkdir /app
WORKDIR /app/

# Add user for running as non-root
RUN addgroup -S nomad-var-dirsync && adduser -S nomad-var-dirsync -G nomad-var-dirsync

# Copy binary in
COPY ./dist/nomad-var-dirsync-TARGETOS-TARGETARCH ./nomad-var-dirsync

# Drop to non-root user
USER nomad-var-dirsync

ENTRYPOINT [ "./nomad-var-dirsync" ]
