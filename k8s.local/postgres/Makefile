.PHONY: scripts
scripts:
	kubectl create configmap postgres-scripts \
		--from-file=./scripts \
		--dry-run=client --output yaml > scripts.yaml
