module github.com/metal3-io/hardware-classification-controller

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/metal3-io/baremetal-operator/apis v0.0.0-20211027194412-408523eff5cc
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime v0.9.7
	sigs.k8s.io/controller-tools v0.4.0
	sigs.k8s.io/kustomize/kustomize/v3 v3.8.5
)
