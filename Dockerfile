FROM alpine:3.4
MAINTAINER Maximilian Pachl <m@ximilian.info>

ADD bin/xkcdbot /usr/sbin/xkcdbot

CMD ["/usr/sbin/xkcdbot"]