.PHONY: config
config:
	kubectl create configmap backend-config \
		--from-file=./config \
		--dry-run=client --output yaml > config.yaml
