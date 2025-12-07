NAMESPACE=eda-workshop

.PHONY: up
up:
	kubectl create namespace $(NAMESPACE) || true
	tilt up --namespace=$(NAMESPACE) --stream=true --host 0.0.0.0

.PHONY: stop
stop:
	tilt down --namespace=$(NAMESPACE)

.PHONY: down
down:
	@echo "🗑️  WARNING: This will DELETE ALL DATA including PVCs!"
	@echo "Press Ctrl+C to cancel or wait 5 seconds to continue..."
	@sleep 5
	tilt down --namespace=$(NAMESPACE) --delete-namespaces
	kubectl delete namespace $(NAMESPACE) || true

.PHONY: setup
setup:
	@bash ./scripts/setup.sh

.PHONY: clean
clean:
	@echo "🧼 Cleaning up Docker system..."
	docker system prune -af --volumes