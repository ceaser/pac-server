FROM ubuntu

ARG PKG=pac-server_*_amd64.deb

RUN echo ${PKG}
COPY ${PKG} .
RUN dpkg -i ${PKG}

EXPOSE 80/tcp
CMD /usr/sbin/pac-server

