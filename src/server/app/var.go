package app

const (
	BuildDev  = "dev"
	BuildProd = "prod"
)

var (
	// Build describes app build type
	Build = BuildDev
	// Version is a semver version of the app
	Version = "0.0.0"
	// CommitHash is latest build commit hash
	CommitHash = "n/a"
	// BuildTimestamp stores when the app was build
	BuildTimestamp = "n/a"
)

// Environment variables
var (
	Port       int
	AccessKey  string
	RefreshKey string
	MongoConn  string
)
