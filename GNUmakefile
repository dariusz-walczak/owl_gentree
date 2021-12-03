


.PHONY: gentree

GENTREE_TAG = owl/gentree:0.1


gentree:
	sudo docker build --tag=${GENTREE_TAG} .
	@echo "==================================================================================================="
	sudo docker run -p 8080:8080/tcp -it --rm ${GENTREE_TAG} --log-level trace
