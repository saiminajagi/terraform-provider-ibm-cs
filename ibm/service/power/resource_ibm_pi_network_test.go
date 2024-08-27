// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
)

func TestAccIBMPINetworkbasic(t *testing.T) {
	name := fmt.Sprintf("tf-pi-network-%d", acctest.RandIntRange(10, 100))
	networkRes := "ibm_pi_network.power_networks"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworkConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists(networkRes),
					resource.TestCheckResourceAttr(networkRes, "pi_network_name", name),
					resource.TestCheckResourceAttrSet(networkRes, "id"),
					resource.TestCheckResourceAttrSet(networkRes, "pi_gateway"),
					resource.TestCheckResourceAttrSet(networkRes, "pi_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPINetworkConfigUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists(networkRes),
					resource.TestCheckResourceAttr(networkRes, "pi_network_name", name),
					resource.TestCheckResourceAttr(networkRes, "pi_dns.#", "1"),
					resource.TestCheckResourceAttrSet(networkRes, "id"),
					resource.TestCheckResourceAttrSet(networkRes, "pi_gateway"),
					resource.TestCheckResourceAttrSet(networkRes, "pi_ipaddress_range.#"),
				),
			},
		},
	})
}
func TestAccIBMPINetworkGatewaybasic(t *testing.T) {
	name := fmt.Sprintf("tf-pi-network-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworkGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_gateway"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPINetworkConfigGatewayUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_gateway", "192.168.17.2"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_ipaddress_range.0.pi_ending_ip_address", "192.168.17.254"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_ipaddress_range.0.pi_starting_ip_address", "192.168.17.3"),
				),
			},
		},
	})
}

func TestAccIBMPINetworkGatewaybasicSatellite(t *testing.T) {
	name := fmt.Sprintf("tf-pi-network-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworkGatewayConfigSatellite(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_gateway"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPINetworkConfigGatewayUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_gateway", "192.168.17.2"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_ipaddress_range.0.pi_ending_ip_address", "192.168.17.254"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_ipaddress_range.0.pi_starting_ip_address", "192.168.17.3"),
				),
			},
		},
	})
}

func TestAccIBMPINetworkDHCPbasic(t *testing.T) {
	name := fmt.Sprintf("tf-pi-network-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworDHCPConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_gateway"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_pi_network.power_networks", "pi_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPINetworkConfigGatewayDHCPUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists("ibm_pi_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_network_name", name),
					resource.TestCheckResourceAttr(
						"ibm_pi_network.power_networks", "pi_dns.#", "1"),
				),
			},
		},
	})
}
func TestAccIBMPINetworkusertags(t *testing.T) {
	name := fmt.Sprintf("tf-pi-network-%d", acctest.RandIntRange(10, 100))
	networkRes := "ibm_pi_network.power_networks"
	networkResData := "data.ibm_pi_network.power_networks_data"
	userTagsString := `["env:dev","test_tag"]`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworkConfigUserTags(name, userTagsString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkExists(networkRes),
					resource.TestCheckResourceAttr(networkResData, "user_tags.#", "2"),
					resource.TestCheckResourceAttr(networkResData, "user_tags.0", "env:test"),
					resource.TestCheckResourceAttr(networkResData, "user_tags.1", "test_tag"),
				),
			},
		},
	})
}

func testAccCheckIBMPINetworkDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_pi_network" {
			continue
		}
		cloudInstanceID, networkID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		networkC := st.NewIBMPINetworkClient(context.Background(), sess, cloudInstanceID)
		_, err = networkC.Get(networkID)
		if err == nil {
			return fmt.Errorf("PI Network still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIBMPINetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPISession()
		if err != nil {
			return err
		}
		cloudInstanceID, networkID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPINetworkClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(networkID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckIBMPINetworkConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id = "%s"
			pi_network_name      = "%s"
			pi_network_type      = "pub-vlan"
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkConfigUpdateDNS(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id = "%s"
			pi_network_name      = "%s"
			pi_network_type      = "pub-vlan"
			pi_dns               = ["127.0.0.1"]
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkGatewayConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id = "%s"
			pi_network_name      = "%s"
			pi_network_type      = "vlan"
			pi_cidr              = "192.168.17.0/24"
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkGatewayConfigSatellite(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id 		= "%s"
			pi_network_name      		= "%s"
			pi_network_type      		= "vlan"
			pi_cidr              		= "192.168.17.0/24"
			pi_network_mtu		 		= 6500
			pi_network_access_config	= "outbound-only"
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkConfigGatewayUpdateDNS(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id = "%s"
			pi_network_name      = "%s"
			pi_network_type      = "vlan"
			pi_dns               = ["127.0.0.1"]
			pi_gateway           = "192.168.17.2"
			pi_cidr              = "192.168.17.0/24"
			pi_ipaddress_range {
				pi_ending_ip_address = "192.168.17.254"
				pi_starting_ip_address = "192.168.17.3"
			}
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworDHCPConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id 		= "%s"
			pi_network_name      		= "%s"
			pi_network_type      		= "dhcp-vlan"
			pi_cidr              		= "10.1.2.0/26"
			pi_dns               		= ["10.1.0.68"]
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkConfigGatewayDHCPUpdateDNS(name string) string {
	return fmt.Sprintf(`
		resource "ibm_pi_network" "power_networks" {
			pi_cloud_instance_id = "%s"
			pi_network_name      = "%s"
			pi_network_type      = "dhcp-vlan"
			pi_cidr              = "10.1.2.0/26"
			pi_dns               = ["10.1.0.69"]
		}
	`, acc.Pi_cloud_instance_id, name)
}

func testAccCheckIBMPINetworkConfigUserTags(name string, userTagsString string) string {
	return fmt.Sprintf(`
		data "ibm_pi_network" "power_networks_data" {
			pi_cloud_instance_id = "%[1]s"
			pi_network_name      = ibm_pi_network.power_networks.pi_network_name
		}

		resource "ibm_pi_network" "power_networks" {
			pi_network_name      = "%[2]s"
			pi_cloud_instance_id = "%[1]s"
			pi_network_type      = "vlan"
			pi_cidr              = "192.168.17.0/24"
			pi_gateway           = "192.168.17.2"
			pi_ipaddress_range {
				pi_ending_ip_address   = "192.168.17.254"
				pi_starting_ip_address = "192.168.17.3"
			}
			pi_user_tags 		 = %[3]s
		}
	`, acc.Pi_cloud_instance_id, name, userTagsString)
}
