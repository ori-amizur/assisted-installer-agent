package dhcp_lease_allocator

import (
	"net"

	"github.com/openshift/assisted-installer-agent/src/util"
	"github.com/openshift/baremetal-runtimecfg/pkg/monitor"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type VIP struct {
	Name       string `yaml:"name"`
	MacAddress string `yaml:"mac-address"`
	IpAddress  string `yaml:"ip-address"`
}

func LeaseVIP(log logrus.FieldLogger, cfgPath, masterDevice, name string, mac net.HardwareAddr, ip string) error {
	iface, err := monitor.LeaseInterface(log, masterDevice, name, mac)

	if err != nil {
		log.WithFields(logrus.Fields{
			"masterDevice": masterDevice,
			"name":         name,
		}).WithError(err).Error("Failed to lease interface")
		return err
	}

	leaseFile := monitor.GetLeaseFile(cfgPath, name)

	// -sf avoiding dhclient from setting the received IP to the interface
	// --no-pid in order to allow running multiple `dhclient` simultaneously
	// -pf allow killing the process
	_, stderr, exitCode := util.Execute("timeout", "5", "dhclient", "-v", "-H", name,
		"-sf", "/bin/true", "-lf", leaseFile,
		"--no-pid", "-1", iface.Name)
	if exitCode != 0 {
		return errors.Errorf("dhclient existed with non-zero exit code %d: %s", exitCode, stderr)
	}
	return nil
}
