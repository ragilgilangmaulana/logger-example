package cservice

const (
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvDevelopment = "development"
	EnvLocal       = "local"
)

const (
	ClusterKubernetes = "kubernetes"
	ClusterLocal      = "local"
)

const (
	NamespaceDefault = "default"
)

const (
	AppDebug        = "APP_DEBUG"
	AppEnv          = "APP_ENV"
	AppKey          = "APP_KEY"
	AppName         = "APP_NAME"
	AppVersion      = "APP_VERSION"
	AppHost         = "APP_HOST"
	AppPort         = "APP_PORT"
	AppEndpoint     = "APP_ENDPOINT"
	AppOpenEndpoint = "APP_OPEN_ENDPOINT"
	AppBasepoint    = "APP_BASEPOINT"
	AppRegistryAddr = "APP_REGISTRY_ADDR"
	AppRegistryPwd  = "APP_REGISTRY_PWD"
	AppTimezone     = "APP_TIMEZONE"
	AppNamespace    = "APP_NAMESPACE"
	AppCluster      = "APP_CLUSTER"

	DBEngine       = "DB_ENGINE"
	DBHost         = "DB_HOST"
	DBPort         = "DB_PORT"
	DBUser         = "DB_USER"
	DBPwd          = "DB_PWD"
	DBName         = "DB_NAME"
	DBSSLMode      = "DB_SSL_MODE"
	DBConnStr      = "DB_CONN_STR"
	DBConnLifetime = "DB_CONN_LIFETIME"
	DBConnMaxIdle  = "DB_CONN_MAX_IDLE"
	DBConnMaxOpen  = "DB_CONN_MAX_OPEN"

	BrokerAddr   = "BROKER_ADDR"
	BrokerKey    = "BROKER_KEY"
	BrokerSecret = "BROKER_SECRET"

	LogTrackLevel   = "LOG_TRACK_LEVEL"
	LogSentryToken  = "LOG_SENTRY_TOKEN"
	LogRollbarToken = "LOG_ROLLBAR_TOKEN"

	APMEnabled                  = "APM_ENABLED"
	APMProviderDatadogAgentHost = "APM_PROVIDER_DATADOG_AGENT_HOST"
	APMProviderDatadogAgentPort = "APM_PROVIDER_DATADOG_AGENT_PORT"

	CacheEnabled            = "CACHE_ENABLED"
	CacheStorePriorities    = "CACHE_STORE_PRIORITIES"
	CacheStoreMemoryEnabled = "CACHE_STORE_MEMORY_ENABLED"
	CacheStoreRedisEnabled  = "CACHE_STORE_REDIS_ENABLED"
	CacheStoreRedisAddr     = "CACHE_STORE_REDIS_ADDR"
	CacheStoreRedisPwd      = "CACHE_STORE_REDIS_PWD"
	CacheStoreRedisDB       = "CACHE_STORE_REDIS_DB"
)
