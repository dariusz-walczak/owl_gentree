


.PHONY: gentree

GENTREE_TAG = owl/gentree:0.1
GENTREE_TEST_TAG = owl/gentree_test:0.1

gentree:
	sudo docker build --tag=${GENTREE_TAG} .
	@echo "==================================================================================================="
	sudo docker run -p 8080:8080/tcp -it --rm ${GENTREE_TAG} --log-level trace


run_ut:
	sudo docker build -f run-unit-tests.Dockerfile --tag=${GENTREE_TEST_TAG} .
	mkdir -p "$(shell pwd)/output"
	sudo docker run --mount type=bind,source="$(shell pwd)/output",target=/output -it --rm ${GENTREE_TEST_TAG} -coverprofile /output/gentree_cover.out
	cd gentree && go tool cover -html=../output/gentree_cover.out -o ../output/cover.html
