package utils

const (
	HomeConfigPath  = "$HOME/.install-tools"
	LocalConfigPath = "./config"

	// DefaultConfigFilename is the default filename of the config file that viper attempts to load (viper.SetConfigName).
	// Viper accepts many suffixes, but only .env was tested. Paths that are searched by viper are defined in ConfigPaths.
	DefaultConfigFilename = "conf"

	// EnvPrefix environment variable sets a prefix of all environment variables bound to command line flags.
	// For example, if the flag is --cloud, the environment variable will be INST_CLOUD.
	EnvPrefix = "INST"
)

var ConfigPaths = []string{
	// First path here has the highest priority.
	HomeConfigPath,
	LocalConfigPath,
}
