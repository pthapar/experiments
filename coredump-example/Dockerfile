FROM alpine:3.9.4
RUN apk add build-base &&  apk add bash
RUN mkdir /app
COPY coredump.c /app/
COPY run.sh /app/
COPY runcrash.sh /app/
WORKDIR /app
RUN gcc coredump.c -o coredump
CMD ["/app/run.sh"]
