package server

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jithin-kg/webpa-common/logging"
	"github.com/jithin-kg/webpa-common/xmetrics"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DefaultPrimaryAddress is the bind address of the primary server (e.g. talaria, petasos, etc)
	DefaultPrimaryAddress = ":8080"

	// DefaultHealthAddress is the bind address of the health check server
	DefaultHealthAddress = ":8081"

	// DefaultMetricsAddress is the bind address of the metrics server
	DefaultMetricsAddress = ":8082"

	// DefaultHealthLogInterval is the interval at which health statistics are emitted
	// when a non-positive log interval is specified
	DefaultHealthLogInterval time.Duration = time.Duration(60 * time.Second)

	// DefaultLogConnectionState is the default setting for logging connection state messages.  This
	// value is primarily used when a *WebPA value is nil.
	DefaultLogConnectionState = false

	// AlternateSuffix is the suffix appended to the server name, along with a period (.), for
	// logging information pertinent to the alternate server.
	AlternateSuffix = "alternate"

	// DefaultProject is used as a metrics namespace when one is not defined.
	DefaultProject = "xmidt"

	// HealthSuffix is the suffix appended to the server name, along with a period (.), for
	// logging information pertinent to the health server.
	HealthSuffix = "health"

	// PprofSuffix is the suffix appended to the server name, along with a period (.), for
	// logging information pertinent to the pprof server.
	PprofSuffix = "pprof"

	// MetricsSuffix is the suffix appended to the server name, along with a period (.), for
	// logging information pertinent to the metrics server.
	MetricsSuffix = "metrics"

	// FileFlagName is the name of the command-line flag for specifying an alternate
	// configuration file for Viper to hunt for.
	FileFlagName = "file"

	// FileFlagShorthand is the command-line shortcut flag for FileFlagName
	FileFlagShorthand = "f"

	// CPUProfileFlagName is the command-line flag for creating a cpuprofile on the server
	CPUProfileFlagName = "cpuprofile"

	// CPUProfileShortHand is the command-line shortcut for creating cpushorthand on the server
	CPUProfileShorthand = "c"

	// MemProfileFlagName is the command-line flag for creating memprofile on the server
	MemProfileFlagName = "memprofile"

	// MemProfileShortHand is the command-line shortcut for creating memprofile on the server
	MemProfileShorthand = "m"
)

// ConfigureFlagSet adds the standard set of WebPA flags to the supplied FlagSet.  Use of this function
// is optional, and necessary only if the standard flags should be supported.  However, this is highly desirable,
// as ConfigureViper can make use of the standard flags to tailor how configuration is loaded or if gathering cpuprofile
// or memprofile data is needed.
func ConfigureFlagSet(applicationName string, f *pflag.FlagSet) {
	// Add the WebPA flags to the FlagSet.
	f.StringP(FileFlagName, FileFlagShorthand, applicationName, "base name of the configuration file")
	f.StringP(CPUProfileFlagName, CPUProfileShorthand, "cpuprofile", "base name of the cpuprofile file")
	f.StringP(MemProfileFlagName, MemProfileShorthand, "memprofile", "base name of the memprofile file")
}

// create CPUProfileFiles creates a cpu profile of the server, its triggered by the optional flag cpuprofile
//
// the CPU profile is created on the server's start
func CreateCPUProfileFile(v *viper.Viper, fp *pflag.FlagSet, l log.Logger) {
	if fp == nil {
		return
	}
	// Look up the cpuprofile flag.
	flag := fp.Lookup("cpuprofile")
	if flag == nil {
		return
	}
	// Create a file to write the profile to.
	f, err := os.Create(flag.Value.String())
	if err != nil {
		l.Log("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		l.Log("could not start CPU profile: ", err)
	}

	defer pprof.StopCPUProfile()
}

// Create CPUProfileFiles creates a memory profile of the server, its triggered by the optional flag memprofile
//
// the memory profile is created on the server's exit.
// this function should be used within the application.
func CreateMemoryProfileFile(v *viper.Viper, fp *pflag.FlagSet, l log.Logger) {
	if fp == nil {
		return
	}

	flag := fp.Lookup("memprofile")
	if flag == nil {
		return
	}

	f, err := os.Create(flag.Value.String())
	if err != nil {
		l.Log("could not create memory profile: ", err)
	}

	defer f.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		l.Log("could not write memory profile: ", err)
	}
}

// ConfigureViper configures a Viper instances using the opinionated WebPA settings.  All WebPA servers should
// use this function.
//
// The flagSet is optional.  If supplied, it will be bound to the given Viper instance.  Additionally, if the
// flagSet has a FileFlagName flag, it will be used as the configuration name to hunt for instead of the
// application name.
/**
ConfigureViper function sets up paths to search for configuration files, binds environment variables to configuration settings,
sets default values for configuration settings,
and optionally binds command-line arguments to configuration settings.
The function takes three arguments: the name of the application, a pointer to a pflag.
FlagSet object containing command-line flags (which can be nil), and a pointer to a viper.
Viper object to be configured. The function returns an error if there is a problem binding command-line arguments to configuration settings.
**/
func ConfigureViper(applicationName string, f *pflag.FlagSet, v *viper.Viper) (err error) {
	// Set up paths to search for configuration files
	v.AddConfigPath(fmt.Sprintf("/etc/%s", applicationName))
	v.AddConfigPath(fmt.Sprintf("$HOME/.%s", applicationName))
	v.AddConfigPath(".")

	// Set up viper to replace dots with underscores in environment variable names, and use the
	// application name as the environment variable prefix
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(applicationName)
	// Set up viper to automatically bind environment variables to configuration settings
	v.AutomaticEnv()
	// Set default values for configuration settings
	v.SetDefault("primary.name", applicationName)
	v.SetDefault("primary.address", DefaultPrimaryAddress)
	v.SetDefault("primary.logConnectionState", DefaultLogConnectionState)

	v.SetDefault("alternate.name", fmt.Sprintf("%s.%s", applicationName, AlternateSuffix))

	v.SetDefault("health.name", fmt.Sprintf("%s.%s", applicationName, HealthSuffix))
	v.SetDefault("health.address", DefaultHealthAddress)
	v.SetDefault("health.logInterval", DefaultHealthLogInterval)
	v.SetDefault("health.logConnectionState", DefaultLogConnectionState)

	v.SetDefault("pprof.name", fmt.Sprintf("%s.%s", applicationName, PprofSuffix))
	v.SetDefault("pprof.logConnectionState", DefaultLogConnectionState)

	v.SetDefault("metric.name", fmt.Sprintf("%s.%s", applicationName, MetricsSuffix))
	v.SetDefault("metric.address", DefaultMetricsAddress)

	v.SetDefault("project", DefaultProject)

	configName := applicationName
	if f != nil {
		if fileFlag := f.Lookup(FileFlagName); fileFlag != nil {
			// use the command-line to specify the base name of the file to be searched for
			configName = fileFlag.Value.String()
		}

		err = v.BindPFlags(f)
	}
	// Set the name of the configuration file to search for (defaults to applicationName)
	v.SetConfigName(configName)
	return
}

/*
Configure is a one-stop shopping function for preparing WebPA configuration.  This function
does not itself read in configuration from the Viper environment.  Typical usage is:

	var (
	  f = pflag.NewFlagSet()
	  v = viper.New()
	)

	if err := server.Configure("petasos", os.Args, f, v); err != nil {
	  // deal with the error, possibly just exiting
	}

	// further customizations to the Viper instance can be done here

	if err := v.ReadInConfig(); err != nil {
	  // more error handling
	}

Usage of this function is only necessary if custom configuration is needed.  Normally,
using New will suffice.
*/
func Configure(applicationName string, arguments []string, f *pflag.FlagSet, v *viper.Viper) (err error) {
	if f != nil {
		ConfigureFlagSet(applicationName, f)
		err = f.Parse(arguments)
		if err != nil {
			return
		}
	}

	err = ConfigureViper(applicationName, f, v)
	return
}

/*
Initialize handles the bootstrapping of the server code for a WebPA node.  It configures Viper,
reads configuration, and unmarshals the appropriate objects.  This function is typically all that's
needed to fully instantiate a WebPA server.  Typical usage:

	var (
	  f = pflag.NewFlagSet()
	  v = viper.New()

	  // can customize both the FlagSet and the Viper before invoking New
	  logger, registry, webPA, err = server.Initialize("petasos", os.Args, f, v)
	)

	if err != nil {
	  // deal with the error, possibly just exiting
	}

Note that the FlagSet is optional but highly encouraged.  If not supplied, then no command-line binding
is done for the unmarshalled configuration.

This function always returns a logger, regardless of any errors.  This allows clients to use the returned
logger when reporting errors.  This function falls back to a logger that writes to os.Stdout if it cannot
create a logger from the Viper environment.
*/

/*
*
Initialize function sets up the Viper environment,
reads in configuration, initializes the logger and metric registry,
and returns a fully initialized WebPA struct along with a logger and metric registry.

Initialize initializes the WebPA server and is the only function required to instantiate the WebPA server.
The function accepts applicationName as a string, arguments as a slice of strings,
f as a pointer to pflag.FlagSet, v as a pointer to viper.Viper, and modules as
*
*/
func Initialize(applicationName string, arguments []string, f *pflag.FlagSet, v *viper.Viper, modules ...xmetrics.Module) (logger log.Logger, registry xmetrics.Registry, webPA *WebPA, err error) {
	// The defer statement is used to ensure that certain statements are executed just before the function returns.
	// The defer statement is used to set the logger variable to a default logger if an error occurs during
	//the execution of the function. This ensures that the function always returns a logger, even if there was an error during initialization.
	//Therefore, using defer in this case is a defensive programming measure to ensure that the caller always has access to a logger,
	//even if the initialization failed.
	defer func() {
		if err != nil {
			// never return a WebPA in the presence of an error, to
			// avoid an ambiguous API
			webPA = nil

			// Make sure there's at least a default logger for the caller to use
			logger = logging.DefaultLogger()
		}
	}()
	// configure viper
	if err = Configure(applicationName, arguments, f, v); err != nil {
		return
	}
	// read in the configuration file
	if err = v.ReadInConfig(); err != nil {
		return
	}

	webPA = &WebPA{
		ApplicationName: applicationName,
	}
	// Here we set rest of proprties of configs to the webpa struct instance, ie  unmarshal configuration settings into webPA struct
	err = v.Unmarshal(webPA)
	if err != nil {
		return
	}
	// create a logger instance
	logger = logging.New(webPA.Log)
	// log a message to indicate that the Viper environment has been initialized
	logger.Log(level.Key(), level.InfoValue(), logging.MessageKey(), "initialized Viper environment", "configurationFile", v.ConfigFileUsed())

	// set the metrics namespace and subsystem to applicationName, if not set in configuration
	// namespace:  In the context of webPA, different servers running in different ports might have different metric
	//and logging configurations, which means they should have different namespaces to differentiate between them.
	//For example, you might have multiple instances of a server running with different configurations, such as different IP addresses, ports,
	//or other settings. Each of these instances might write metrics to the same Prometheus server, but they should use different namespaces to
	//differentiate between them. This allows you to easily identify which instance is having issues, or which instance is generating the most traffic.
	if len(webPA.Metric.MetricsOptions.Namespace) == 0 {
		webPA.Metric.MetricsOptions.Namespace = applicationName
	}

	if len(webPA.Metric.MetricsOptions.Subsystem) == 0 {
		webPA.Metric.MetricsOptions.Subsystem = applicationName
	}
	// set the logger for the metrics registry to the created logger instance
	webPA.Metric.MetricsOptions.Logger = logger
	// create the metrics registry
	registry, err = webPA.Metric.NewRegistry(modules...)
	if err != nil {
		return
	}
	// create a CPU profile file, if configured
	CreateCPUProfileFile(v, f, logger)

	return
}
