FROM alpine:3.4
MAINTAINER Maximilian Pachl <m@ximilian.info>

RUN apk --update add ca-certificates

ADD bin/xkcdbot /usr/sbin/xkcdbot
RUN chmod 755 /usr/sbin/xkcdbot

CMD ["/usr/sbin/xkcdbot"]
