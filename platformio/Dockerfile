FROM debian:8.6
MAINTAINER Matthias Pfeil <matthias.pfeil@wwu.de>

RUN apt-get update && apt-get install -y curl python2.7 && rm -rf /var/lib/apt/lists/*
RUN python2.7 -c "$(curl -fsSL https://raw.githubusercontent.com/platformio/platformio/master/scripts/get-platformio.py)"

ENV PLATFORMIO_ENVS_DIR /sensebox/build
ENV BOARD uno
ENV SENSEBOX_RESOURCES_VERSION 1.0.0

WORKDIR sensebox

RUN pio platforms install atmelavr && pio init --board $BOARD

RUN curl -L https://github.com/sensebox/resources/archive/v$SENSEBOX_RESOURCES_VERSION.tar.gz | tar zxvf - resources-$SENSEBOX_RESOURCES_VERSION/libraries && \ 
	mv resources-$SENSEBOX_RESOURCES_VERSION/libraries/*/ lib/ && \ 
	rm -rf resources-$SENSEBOX_RESOURCES_VERSION/

CMD ["pio","run"]
