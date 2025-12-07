.PHONY: config
config:
	kubectl create configmap seaweed-config \
		--from-file=./config \
		--dry-run=client --output yaml > config.yaml
