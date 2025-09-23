package config

// rest api client envs.
const (
	CalendarRestAPIHostEnv     = "CALENDAR_REST_API_HOST"
	CalendarRestAPIHostDefault = "localhost"
	CalendarRestAPIPortEnv     = "CALENDAR_REST_API_PORT"
	CalendarRestAPIPortDefault = "8888"
)

// grpc api client  envs.
const (
	CalendarGrpcAPIHostEnv     = "CALENDAR_GRPC_API_HOST"
	CalendarGrpcAPIHostDefault = "localhost"
	CalendarGrpcAPIPortEnv     = "CALENDAR_GRPC_API_PORT"
	CalendarGrpcAPIPortDefault = "50051"
)

// repo envs.
const (
	CalendarDBHostEnv         = "CALENDAR_DB_HOST"
	CalendarDBHostDefault     = "localhost"
	CalendarDBPortEnv         = "CALENDAR_DB_PORT"
	CalendarDBPortDefault     = "5432"
	CalendarDBNameEnv         = "CALENDAR_DB_NAME"
	CalendarDBNameDefault     = "calendar"
	CalendarDBUserEnv         = "CALENDAR_DB_USER"
	CalendarDBUserDefault     = "postgres"
	CalendarDBPasswordEnv     = "CALENDAR_DB_PASSWORD"
	CalendarDBPasswordDefault = "password"
)
