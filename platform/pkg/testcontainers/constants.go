package testcontainers

const (
	MongoContainerName = "inventory-mongo-test"
	MongoPort          = "27017"

	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_INITDB_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"

	// Logger
	LoggerLevelKey  = "LOGGER_LEVEL"
	LoggerAsJsonKey = "LOGGER_AS_JSON"

	// gRPC
	GrpcPortKey = "GRPC_PORT"
)
