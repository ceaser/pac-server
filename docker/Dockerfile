FROM ubuntu

ARG PKG=pac-server_*_amd64.deb

RUN echo ${PKG}
COPY ${PKG} .
RUN dpkg -i ${PKG}

EXPOSE 80/tcp
ENTRYPOINT /usr/sbin/pac-server

