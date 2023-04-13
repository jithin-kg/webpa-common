package servicecfg

import (
	"github.com/jithin-kg/webpa-common/service/consul"
	"github.com/jithin-kg/webpa-common/service/zk"
)

// resetEnvironmentFactories resets the global factories for service.Environment objects
func resetEnvironmentFactories() {
	zookeeperEnvironmentFactory = zk.NewEnvironment
	consulEnvironmentFactory = consul.NewEnvironment
}
