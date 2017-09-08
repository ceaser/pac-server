FROM ubuntu

COPY pac-server_0.0.1-8_amd64.deb .
RUN dpkg -i pac-server_0.0.1-8_amd64.deb

EXPOSE 80/tcp
CMD /usr/sbin/pac-server

