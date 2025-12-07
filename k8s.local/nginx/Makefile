.PHONY: config
config:
	kubectl create configmap nginx-config \
		--from-file=./config \
		--dry-run=client --output yaml > config.yaml
