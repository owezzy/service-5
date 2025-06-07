# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"


# ==============================================================================
# Define dependencies

GOLANG          := golang:1.22
ALPINE          := alpine:3.19
KIND            := kindest/node:v1.32.2
POSTGRES        := postgres:15.4
VAULT           := hashicorp/vault:1.15
GRAFANA         := grafana/grafana:10.1.0
PROMETHEUS      := prom/prometheus:v2.47.0
TEMPO           := grafana/tempo:2.2.0
LOKI            := grafana/loki:2.9.0
PROMTAIL        := grafana/promtail:2.9.0

KIND_CLUSTER    := ardan-starter-cluster
NAMESPACE       := sales-system
APP             := sales
BASE_IMAGE_NAME := ardanlabs/service
SERVICE_NAME    := sales-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME)-metrics:$(VERSION)


# ==============================================================================
# dependencies setup

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew tap hashicorp/tap
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list vault || brew install hashicorp/tap/vault

dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(KIND)
	docker pull $(POSTGRES)
	docker pull $(VAULT)
	docker pull $(GRAFANA)
	docker pull $(PROMETHEUS)
	docker pull $(TEMPO)
	docker pull $(LOKI)
	docker pull $(PROMTAIL)


# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view:
	expvarmon -ports="localhost:3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"



# ==============================================================================
# Building containers

all: service

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.




# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)
	kind load docker-image $(GRAFANA) --name $(KIND_CLUSTER)
	kind load docker-image $(PROMETHEUS) --name $(KIND_CLUSTER)
	kind load docker-image $(TEMPO) --name $(KIND_CLUSTER)
	kind load docker-image $(LOKI) --name $(KIND_CLUSTER)
	kind load docker-image $(PROMTAIL) --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	cd zarf/k8s/dev/sales; kustomize edit set image service-image=$(SERVICE_IMAGE)
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

	cd zarf/k8s/dev/sales; kustomize edit set image metrics-image=$(METRICS_IMAGE)
	kind load docker-image $(METRICS_IMAGE) --name $(KIND_CLUSTER)
# podman is currently experimental, and fails for some reason with kind load
# docker-image (possibly a tagging issue?) but the below works.

dev-apply:
	kustomize build zarf/k8s/dev/vault | kubectl apply -f -

	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/grafana | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=grafana --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/prometheus | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=prometheus --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/tempo | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=tempo --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/loki | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=loki --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/promtail | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=promtail --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

# run on code change
dev-update: all dev-load dev-restart

# run on YAML config change
dev-update-apply:  all dev-load dev-apply

# ------------------------------------------------------------------------------
dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)


dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100


pgcli:
	pgcli postgresql://postgres:postgres@localhost

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-clean-cache:
	go clean -modcache

list:
	go list -mod=mod all


# ==============================================================================
# Class Stuff

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

run-help:
	go run app/services/sales-api/main.go --help | go run app/tooling/logfmt/main.go

curl:
	curl -il http://localhost:3000/v1/hack

curl-auth:
	curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/hackauth

load-hack:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/hack"

admin:
	go run app/tooling/sales-admin/main.go

ready:
	curl -il http://localhost:3000/v1/readiness

live:
	curl -il http://localhost:3000/v1/liveness

curl-create:
	curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"bill","email":"b@gmail.com","roles":["ADMIN"],"department":"IT","password":"123","passwordConfirm":"123"}' http://localhost:3000/v1/users


# ==============================================================================
# Admin Frontend

ADMIN_FRONTEND_PREFIX := ./app/frontends/admin

write-token-to-env:
	echo "VITE_SERVICE_API=http://localhost:3000/v1" > ${ADMIN_FRONTEND_PREFIX}/.env
	make token | grep -o '"ey.*"' | awk '{print "VITE_SERVICE_TOKEN="$$1}' >> ${ADMIN_FRONTEND_PREFIX}/.env

admin-gui-install:
	pnpm -C ${ADMIN_FRONTEND_PREFIX} install

admin-gui-dev: admin-gui-install
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run dev

admin-gui-build: admin-gui-install
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run build

admin-gui-start-build: admin-gui-build
	pnpm -C ${ADMIN_FRONTEND_PREFIX} run preview

admin-gui-run: write-token-to-env admin-gui-start-build

# ==============================================================================
# Running using Service Weaver.

wea-dev-gotooling: dev-gotooling
	go install github.com/ServiceWeaver/weaver/cmd/weaver@latest
	go install github.com/ServiceWeaver/weaver-kube/cmd/weaver-kube@latest

wea-dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl --context=kind-$(KIND_CLUSTER) wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

wea-dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

wea-dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl --context=kind-$(KIND_CLUSTER) apply -f -
	kubectl rollout status --context=kind-$(KIND_CLUSTER) --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	cd app/weaver/sales-api; GOOS=linux GOARCH=amd64 go build .
	$(eval WEAVER_YAML := $(shell weaver-kube deploy app/weaver/sales-api/dev.toml))
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

	kubectl --context=kind-$(KIND_CLUSTER) apply -f $(WEAVER_YAML)
	kubectl wait pods --namespace=$(NAMESPACE) --selector appName=$(APP)-api --timeout=120s --for=condition=Ready
