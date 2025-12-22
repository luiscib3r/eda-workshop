.PHONY: config
config:
	kubectl create configmap jaeger-config \
		--from-file=config.yaml=jaeger.yaml \
		--dry-run=client --output yaml > jaeger-config.yaml
