.PHONY: config
config:
	kubectl create configmap ocr-llm-config \
		--from-file=./config \
		--dry-run=client --output yaml > config.yaml
