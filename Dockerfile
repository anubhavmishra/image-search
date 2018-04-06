FROM alpine:3.1
MAINTAINER Anubhav Mishra <anubhavmishra@me.com>
ADD build/linux/amd64/image-search /usr/bin/image-search
ENTRYPOINT ["image-search"]