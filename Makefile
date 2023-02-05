
IMAGE_TAG = peary:stable
CONTAINER_NAME = peary_container

build:
	docker build --tag=peary:stable --rm .

run:
	-docker rm ${CONTAINER_NAME}

	docker run \
    	--name=${CONTAINER_NAME} \
    	--env-file=.env \
    	-v /home/$(shell id -nu 1000)/volumes/peary_data:/data \
    	-p 8080:8080 \
    ${IMAGE_TAG}

clear:
	-docker rm -f ${CONTAINER_NAME}
	-docker image rm -f ${IMAGE_TAG}
