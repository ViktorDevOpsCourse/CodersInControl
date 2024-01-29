
help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                              \
		nb = sub( /^## /, "", helpMsg );                      \
		if(nb == 0) {                                         \
			helpMsg = $$0;                                    \
			nb = sub( /^[^:]+:.* ## /, "", helpMsg );         \
		}                                                     \
		if (nb)                                               \
			printf "\033[36m%-20s\033[0m %s\n", $$1, helpMsg; \
	}                                                         \
	{ helpMsg = $$0 }'                                        \
	$(MAKEFILE_LIST)

start: ## Create KinD clusters and deploy Flux
	./scripts/start.sh

test: ## Run tests of installed Podinfo applications on different KinD clusters
	./scripts/test.sh

clean: ## Remove KinD clusters and clean Flux manifests
	./scripts/clean.sh
