FROM ubuntu:16.04
MAINTAINER Dmitrii Beliakov

RUN mkdir /srv/app
COPY api/api /srv/app/revisor
COPY client/dist /srv/app/client

EXPOSE 80

WORKDIR /srv/app
CMD ["/srv/app/revisor"]
