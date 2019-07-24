ARG NODE_VERSION=8.15.0-slim
FROM node:${NODE_VERSION}
ARG GITHUB_ACCOUNT=telefonicaid
ARG GITHUB_REPOSITORY=iotagent-json
ARG DOWNLOAD=latest
ARG SOURCE_BRANCH=master

# Copying Build time arguments to environment variables so they are persisted at run time and can be 
# inspected within a running container.
# see: https://vsupalov.com/docker-build-time-env-values/  for a deeper explanation.

ENV GITHUB_ACCOUNT=${GITHUB_ACCOUNT}
ENV GITHUB_REPOSITORY=${GITHUB_REPOSITORY}
ENV DOWNLOAD=${DOWNLOAD}

MAINTAINER FIWARE IoTAgent Team. Telefónica I+D

# IMPORTANT: For a Production Environment Use Docker Secrets to define
#  these values and add _FILE to the name of the variable.
#
# - IOTA_AUTH_USER, IOTA_AUTH_PASSWORD - when using Keystone Security 
# - IOTA_AUTH_CLIENT_ID, IOTA_AUTH_CLIENT_SECRET - when using OAuth2 Security

# The following RUN command retrieves the source code from GitHub.
# 
# To obtain the latest stable release run this Docker file with the parameters
# --no-cache --build-arg DOWNLOAD=stable
# To obtain any speciifc version of a release run this Docker file with the parameters
# --no-cache --build-arg DOWNLOAD=1.7.0
#
# The default download is the latest tip of the master of the named repository on GitHub
#
# Alternatively for local development, just copy this Dockerfile into file the root of the repository and 
# replace the whole RUN statement by the following COPY statement in your local source using :
#
# COPY . /opt/iotajson/
#
RUN if [ "${DOWNLOAD}" = "latest" ] ; \
	then \
		RELEASE="${SOURCE_BRANCH}"; \		
		echo "INFO: Building Latest Development from ${SOURCE_BRANCH} branch."; \
	elif [ "${DOWNLOAD}" = "stable" ]; \
	then \
		RELEASE=$(curl -s https://api.github.com/repos/"${GITHUB_ACCOUNT}"/"${GITHUB_REPOSITORY}"/releases/latest | grep 'tag_name' | cut -d\" -f4); \
		echo "INFO: Building Latest Stable Release: ${RELEASE}"; \
	else \
	 	RELEASE="${DOWNLOAD}"; \
	 	echo "INFO: Building Release: ${RELEASE}"; \
	fi && \
	RELEASE_CONCAT=$(echo "$SOURCE_BRANCH" | tr / -); \
	# Ensure that unzip is installed, and download the sources
	apt-get update && \
	apt-get install -y  --no-install-recommends unzip && \
	wget --no-check-certificate -O source.zip https://github.com/"${GITHUB_ACCOUNT}"/"${GITHUB_REPOSITORY}"/archive/"${RELEASE}".zip && \
	unzip source.zip && \
	rm source.zip && \
	mv "${GITHUB_REPOSITORY}-${RELEASE_CONCAT}" /opt/iotajson && \
	# Remove unzip and clean apt cache
	apt-get clean && \
	apt-get remove -y unzip && \
	apt-get -y autoremove
	
WORKDIR /opt/iotajson

RUN \
	# Ensure that Git is installed prior to running npm install
	apt-get update -y && \
	apt-get install -y git && \
	npm install pm2@3.2.2 -g && \
	echo "INFO: npm install --production..." && \
	npm install --production && \
        npm install axios && \
        # Install node.js
        apt-get install curl -y && \
        apt-get install software-properties-common -y && \
        curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
        apt-get install nodejs -y && \
        node -v && \
	# Remove Git and clean apt cache
	apt-get clean && \
	apt-get remove -y git && \
	apt-get -y autoremove && \ 
	chmod +x docker/entrypoint.sh
        
COPY ./function.js ./function.js
COPY ./main.js ./main.js
COPY ./iota-config.sh ./iota-config.sh
COPY ./ngsi ./ngsi

RUN chmod +x ./iota-config.sh && \
    chmod 777 ./config.js

USER root
ENV NODE_ENV=production

# Expose Ports - Default to 4041 for NORTH PORT 4041, 7896 for HTTP PORT
EXPOSE ${IOTA_NORTH_PORT:-4041} ${IOTA_HTTP_PORT:-7896}

ENTRYPOINT [ "node","main.js" ]
