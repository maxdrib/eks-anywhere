package hardware_test

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/aws/eks-anywhere/pkg/providers/tinkerbell/hardware"
)

func TestValidMachine(t *testing.T) {
	g := gomega.NewWithT(t)

	machine := hardware.Machine{
		Id:           "unique string",
		IpAddress:    "10.10.10.10",
		Gateway:      "10.10.10.1",
		Nameservers:  []string{"nameserver1"},
		Netmask:      "255.255.255.255",
		MacAddress:   "00:00:00:00:00:00",
		Hostname:     "localhost",
		BmcIpAddress: "10.10.10.11",
		BmcUsername:  "username",
		BmcPassword:  "password",
		BmcVendor:    "dell",
	}

	err := machine.Validate()

	g.Expect(err).ToNot(gomega.HaveOccurred())
}

func TestInvalidMachine(t *testing.T) {
	g := gomega.NewWithT(t)

	cases := map[string]hardware.Machine{
		"EmptyID": {
			Id:           "",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyIpAddress": {
			Id:           "unique string",
			IpAddress:    "",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"InvalidIpAddress": {
			Id:           "unique string",
			IpAddress:    "invalid",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyGateway": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"InvalidGateway": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "invalid",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"NoNameservers": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyNameserver": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{""},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyNetmask": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyMacAddress": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"InvalidMacAddress": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "invalid mac",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyHostname": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyBmcIpAddress": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"InvalidBmcIpAddress": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "invalid",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyBmcUsername": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "",
			BmcPassword:  "password",
			BmcVendor:    "dell",
		},
		"EmptyBmcPassword": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "",
			BmcVendor:    "dell",
		},
		"EmptyBmcVendor": {
			Id:           "unique string",
			IpAddress:    "10.10.10.10",
			Gateway:      "10.10.10.1",
			Nameservers:  []string{"nameserver1"},
			Netmask:      "255.255.255.255",
			MacAddress:   "00:00:00:00:00:00",
			Hostname:     "localhost",
			BmcIpAddress: "10.10.10.11",
			BmcUsername:  "username",
			BmcPassword:  "password",
			BmcVendor:    "",
		},
	}

	for name, machine := range cases {
		t.Run(name, func(t *testing.T) {
			err := machine.Validate()
			g.Expect(err).To(gomega.HaveOccurred())
		})
	}
}