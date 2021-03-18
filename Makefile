
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:crdVersions=v1,trivialVersions=true,preserveUnknownFields=false"

BIN_DIR := $(PWD)/tools/bin
KUSTOMIZE := $(BIN_DIR)/kustomize
CONTROLLER_GEN := $(BIN_DIR)/controller-gen
COVER_PROFILE = cover.out

all: manager

# Run tests
test: generate fmt vet manifests unit

unit:
	go test ./... -coverprofile $(COVER_PROFILE)

unit-profile: unit
	go tool cover -func=$(COVER_PROFILE)

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: $(KUSTOMIZE) manifests
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: $(KUSTOMIZE) manifests
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: $(KUSTOMIZE) manifests
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

$(KUSTOMIZE):
	./tools/install_kustomize.sh

# find or download controller-gen
# download controller-gen if necessary
$(CONTROLLER_GEN):
	./tools/install_controller_gen.sh
