package cservice

type ServiceClass string

const (
	ServiceClassSuper   ServiceClass = "super-service"
	ServiceClassCore    ServiceClass = "core-service"
	ServiceClassProduct ServiceClass = "product-service"
	ServiceClassUtility ServiceClass = "utility-service"
)

type ServiceMetaKey string

const (
	ServiceMetaKeyEnv         ServiceMetaKey = "svc-env"
	ServiceMetaKeyMachineName ServiceMetaKey = "svc-machine-name"
	ServiceMetaMainDataEngine ServiceMetaKey = "svc-main-data-engine"
)
