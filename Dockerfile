FROM quay.io/podman/stable

COPY . /app
WORKDIR /app

#RUN chgrp -R 0 /app && \
#    chmod -R g=u /app \

RUN ./install-deps.sh && make build
