FROM alpine

#RUN apk update && \
#    apk add --no-cache tzdata

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY ./bld/upnp-exporter-$TARGETOS-$TARGETARCH /app/upnp-exporter

EXPOSE 8080
ENTRYPOINT ["/app/upnp-exporter"]
CMD []
