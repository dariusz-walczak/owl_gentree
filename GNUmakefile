


.PHONY: gentree

GENTREE_TAG = owl/gentree:0.1
GENTREE_TEST_TAG = owl/gentree_test:0.1
GENTREE_LINT_TAG = owl/gentree_lint:0.1

gentree:
	sudo docker build --tag=${GENTREE_TAG} .
	@echo "==================================================================================================="
	sudo docker run -p 8080:8080/tcp -it --rm ${GENTREE_TAG} --log-level trace

OUTPUT_PATH = "$(shell pwd)/output"

run_ut:
	sudo docker build -f run-unit-tests.Dockerfile --tag=${GENTREE_TEST_TAG} .
	mkdir -p ${OUTPUT_PATH}
	sudo docker run --mount type=bind,source="$(shell pwd)/output",target=/output -it --rm \
		${GENTREE_TEST_TAG}
	sudo chown --reference=${OUTPUT_PATH} \
		${OUTPUT_PATH}/cover.html ${OUTPUT_PATH}/gentree_cover.out

run_linter:
	sudo docker build -f run-linter.Dockerfile --tag=${GENTREE_LINT_TAG} .
	sudo docker run -it --rm ${GENTREE_LINT_TAG}
