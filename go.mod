module github.com/metal3-io/hardware-classification-controller

go 1.16

require (
	github.com/go-logr/logr v0.2.1
	github.com/metal3-io/baremetal-operator v0.0.0-20201006073612-56a49dc7016a
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
	k8s.io/klog v1.0.0
	sigs.k8s.io/cluster-api v0.3.6
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/controller-tools v0.4.0
	sigs.k8s.io/kustomize/kustomize/v3 v3.8.5
)
