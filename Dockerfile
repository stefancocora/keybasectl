FROM alpine:3.8

MAINTAINER "Stefan Cocora <stefan.cocora@googlemail.com>"

ENV ELF_NAME="keybasectl"
ENV UNPRIVILEDGED_GROUP=users
ENV UNPRIVILEDGED_USER="keybasectl"
ENV HOME=/home/keybasectl
ENV CWD="/${HOME}"

# pkgs
RUN apk add --no-cache bash git

# add non-root user
RUN adduser -u 1000 -G ${UNPRIVILEDGED_GROUP}  -s /bin/sh -D -h /home/${UNPRIVILEDGED_USER} ${UNPRIVILEDGED_USER}
RUN addgroup ${UNPRIVILEDGED_USER} ${UNPRIVILEDGED_GROUP}

COPY Dockerfile /Dockerfile

COPY ./bin/go/${ELF_NAME}* /usr/bin/keybasectl
RUN chmod 755 /usr/bin/${ELF_NAME}

RUN chown ${UNPRIVILEDGED_USER}:${UNPRIVILEDGED_GROUP} -R ${HOME}
USER ${UNPRIVILEDGED_USER}

WORKDIR ${CWD}

# ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["keybasectl", "-h"]
