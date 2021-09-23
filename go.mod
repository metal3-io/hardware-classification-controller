module github.com/metal3-io/hardware-classification-controller

go 1.16

require (
	github.com/go-logr/logr v0.2.1
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/metal3-io/baremetal-operator v0.0.0-20201006073612-56a49dc7016a
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.4.0
	sigs.k8s.io/kustomize/kustomize/v3 v3.8.5
)
