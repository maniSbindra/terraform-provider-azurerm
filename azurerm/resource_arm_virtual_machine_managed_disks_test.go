package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_explicit(t *testing.T) {
	var vm compute.VirtualMachine
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_explicit(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_implicit(t *testing.T) {
	var vm compute.VirtualMachine
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_implicit(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_attach(t *testing.T) {
	var vm compute.VirtualMachine
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_attach(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_withDataDisk_managedDisk_explicit(t *testing.T) {
	var vm compute.VirtualMachine

	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_withDataDisk_managedDisk_explicit(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_withDataDisk_managedDisk_implicit(t *testing.T) {
	var vm compute.VirtualMachine

	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_withDataDisk_managedDisk_implicit(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_deleteManagedDiskOptOut(t *testing.T) {
	var vm compute.VirtualMachine
	var osd string
	var dtd string
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMVirtualMachine_withDataDisk_managedDisk_implicit(ri, location)
	postConfig := testAccAzureRMVirtualMachine_basicLinuxMachineDeleteVM_managedDisk(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
					testLookupAzureRMVirtualMachineManagedDiskID(&vm, "myosdisk1", &osd),
					testLookupAzureRMVirtualMachineManagedDiskID(&vm, "mydatadisk1", &dtd),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineManagedDiskExists(&osd, true),
					testCheckAzureRMVirtualMachineManagedDiskExists(&dtd, true),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_deleteManagedDiskOptIn(t *testing.T) {
	var vm compute.VirtualMachine
	var osd string
	var dtd string
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_DestroyDisksBefore(ri, location)
	postConfig := testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_DestroyDisksAfter(ri, location)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineExists("azurerm_virtual_machine.test", &vm),
					testLookupAzureRMVirtualMachineManagedDiskID(&vm, "myosdisk1", &osd),
					testLookupAzureRMVirtualMachineManagedDiskID(&vm, "mydatadisk1", &dtd),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualMachineManagedDiskExists(&osd, false),
					testCheckAzureRMVirtualMachineManagedDiskExists(&dtd, false),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_osDiskTypeConflict(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_osDiskTypeConflict(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Conflict between `vhd_uri`"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_dataDiskTypeConflict(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualMachine_dataDiskTypeConflict(ri, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Conflict between `vhd_uri`"),
			},
		},
	})
}

func TestAccAzureRMVirtualMachine_bugAzureRM33(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(7)
	config := testAccAzureRMVirtualMachine_bugAzureRM33(ri, rs, testLocation())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}

func testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_explicit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        disk_size_gb = "50"
        managed_disk_type = "Standard_LRS"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_implicit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        disk_size_gb = "50"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_attach(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_managed_disk" "test" {
    name = "acctmd-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    storage_account_type = "Standard_LRS"
    create_option = "Empty"
    disk_size_gb = "1"
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        disk_size_gb = "50"
        managed_disk_type = "Standard_LRS"
    }

    storage_data_disk {
        name = "${azurerm_managed_disk.test.name}"
    	create_option = "Attach"
    	disk_size_gb = "1"
    	lun = 0
        managed_disk_id = "${azurerm_managed_disk.test.id}"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_DestroyDisksBefore(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    delete_os_disk_on_termination = true

    storage_data_disk {
        name          = "mydatadisk1"
    	disk_size_gb  = "1"
    	create_option = "Empty"
    	lun           = 0
    }

    delete_data_disks_on_termination = true

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_basicLinuxMachine_managedDisk_DestroyDisksAfter(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_basicLinuxMachineDeleteVM_managedDisk(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_withDataDisk_managedDisk_explicit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    storage_data_disk {
        name          = "dtd-%d"
    	disk_size_gb  = "1"
    	create_option = "Empty"
        caching       = "ReadWrite"
    	lun           = 0
    	managed_disk_type = "Standard_LRS"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_withDataDisk_managedDisk_implicit(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "myosdisk1"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    storage_data_disk {
        name          = "mydatadisk1"
    	disk_size_gb  = "1"
    	create_option = "Empty"
        caching       = "ReadWrite"
    	lun           = 0
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_osDiskTypeConflict(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        disk_size_gb = "10"
        managed_disk_type = "Standard_LRS"
        vhd_uri = "should_cause_conflict"
    }

    storage_data_disk {
        name = "mydatadisk1"
        caching = "ReadWrite"
        create_option = "Empty"
        disk_size_gb = "45"
        managed_disk_type = "Standard_LRS"
        lun = "0"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_dataDiskTypeConflict(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_D1_v2"

    storage_image_reference {
	publisher = "Canonical"
	offer = "UbuntuServer"
	sku = "14.04.2-LTS"
	version = "latest"
    }

    storage_os_disk {
        name = "osd-%d"
        caching = "ReadWrite"
        create_option = "FromImage"
        disk_size_gb = "10"
        managed_disk_type = "Standard_LRS"
    }

    storage_data_disk {
        name = "mydatadisk1"
        caching = "ReadWrite"
        create_option = "Empty"
        disk_size_gb = "45"
        managed_disk_type = "Standard_LRS"
        lun = "0"
    }

    storage_data_disk {
        name = "mydatadisk1"
        vhd_uri = "should_cause_conflict"
        caching = "ReadWrite"
        create_option = "Empty"
        disk_size_gb = "45"
        managed_disk_type = "Standard_LRS"
        lun = "1"
    }

    os_profile {
	computer_name = "hn%d"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_linux_config {
	disable_password_authentication = false
    }

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMVirtualMachine_bugAzureRM33(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctvn-%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctsub-%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
    name = "acctni-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    ip_configuration {
    	name = "testconfiguration1"
    	subnet_id = "${azurerm_subnet.test.id}"
    	private_ip_address_allocation = "dynamic"
    }
}

resource "azurerm_virtual_machine" "test" {
    name = "acctvm%s"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    network_interface_ids = ["${azurerm_network_interface.test.id}"]
    vm_size = "Standard_F1"

    storage_image_reference {
      publisher = "MicrosoftWindowsServer"
      offer     = "WindowsServer"
      sku       = "2012-Datacenter"
      version   = "latest"
    }

    storage_os_disk {
      name              = "myosdisk1"
      caching           = "ReadWrite"
      create_option     = "FromImage"
      managed_disk_type = "Standard_LRS"
    }

    os_profile {
	computer_name = "acctvm%s"
	admin_username = "testadmin"
	admin_password = "Password1234!"
    }

    os_profile_windows_config {}

    tags {
    	environment = "Production"
    	cost-center = "Ops"
    }
}
`, rInt, location, rInt, rInt, rInt, rString, rString)
}

func testCheckAzureRMVirtualMachineManagedDiskExists(managedDiskID *string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		d, err := testGetAzureRMVirtualMachineManagedDisk(managedDiskID)
		if err != nil {
			return fmt.Errorf("Error trying to retrieve Managed Disk %s, %+v", *managedDiskID, err)
		}
		if d.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Unable to find Managed Disk %s", *managedDiskID)
		}
		if d.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Found unexpected Managed Disk %s", *managedDiskID)
		}

		return nil
	}
}

func testLookupAzureRMVirtualMachineManagedDiskID(vm *compute.VirtualMachine, diskName string, managedDiskID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if osd := vm.StorageProfile.OsDisk; osd != nil {
			if strings.EqualFold(*osd.Name, diskName) {
				if osd.ManagedDisk != nil {
					id, err := findAzureRMVirtualMachineManagedDiskID(osd.ManagedDisk)
					if err != nil {
						return fmt.Errorf("Unable to parse Managed Disk ID for OS Disk %s, %+v", diskName, err)
					}
					*managedDiskID = id
					return nil
				}
			}
		}

		for _, dataDisk := range *vm.StorageProfile.DataDisks {
			if strings.EqualFold(*dataDisk.Name, diskName) {
				if dataDisk.ManagedDisk != nil {
					id, err := findAzureRMVirtualMachineManagedDiskID(dataDisk.ManagedDisk)
					if err != nil {
						return fmt.Errorf("Unable to parse Managed Disk ID for Data Disk %s, %+v", diskName, err)
					}
					*managedDiskID = id
					return nil
				}
			}
		}

		return fmt.Errorf("Unable to locate disk %s on vm %s", diskName, *vm.Name)
	}
}

func findAzureRMVirtualMachineManagedDiskID(md *compute.ManagedDiskParameters) (string, error) {
	_, err := parseAzureResourceID(*md.ID)
	if err != nil {
		return "", err
	}
	return *md.ID, nil
}

func testGetAzureRMVirtualMachineManagedDisk(managedDiskID *string) (*disk.Model, error) {
	armID, err := parseAzureResourceID(*managedDiskID)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Managed Disk ID %s, %+v", *managedDiskID, err)
	}
	name := armID.Path["disks"]
	resourceGroup := armID.ResourceGroup
	conn := testAccProvider.Meta().(*ArmClient).diskClient
	d, err := conn.Get(resourceGroup, name)
	//check status first since sdk client returns error if not 200
	if d.Response.StatusCode == http.StatusNotFound {
		return &d, nil
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}
