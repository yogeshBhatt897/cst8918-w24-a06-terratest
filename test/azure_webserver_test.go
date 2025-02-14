package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Subscription ID (Replace with your actual Subscription ID)
var subscriptionID string = "59020956-debc-48d5-a0ba-f59460247e75"

func TestAzureLinuxVMCreation(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "bhat0199",
		},
	}

	// Run Terraform Init and Apply
	terraform.InitAndApply(t, terraformOptions)

	// Ensure Terraform is destroyed after tests
	defer terraform.Destroy(t, terraformOptions)

	// Retrieve Terraform output values dynamically
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	publicIP := terraform.Output(t, terraformOptions, "public_ip")

	// ✅ Wait for VM resources to stabilize
	time.Sleep(10 * time.Second)

	// ✅ Test 1: Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, subscriptionID, resourceGroupName, vmName), "VM does not exist")

	// ✅ Test 2: Validate the VM is Running Ubuntu
	vm := azure.GetVirtualMachine(t, subscriptionID, resourceGroupName, vmName)
	require.NotNil(t, vm, "Failed to fetch VM details")

	osProfile := vm.StorageProfile.ImageReference
	assert.Equal(t, "Canonical", *osProfile.Publisher, "OS Publisher mismatch")
	assert.Equal(t, "UbuntuServer", *osProfile.Offer, "OS Offer mismatch")

	// ✅ Test 3: Check if NIC is Connected to VM
	// ✅ Test 3: Check if NIC is Connected to VM
nic, err := azure.GetNetworkInterfaceE(subscriptionID, resourceGroupName, nicName)
require.NoError(t, err, "Failed to fetch NIC details")
require.NotNil(t, nic.VirtualMachine, "NIC is not attached to a VM")


	// ✅ Test 4: Validate Public IP Address is Assigned
	require.NotEmpty(t, publicIP, "Public IP is not assigned")
}
