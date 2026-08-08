package driver

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
)

// baseValidFlags returns the minimal set of flag values required for
// SetConfigFromFlags to succeed, so tests only need to override what they
// care about.
func baseValidFlags(overrides map[string]interface{}) *drivers.CheckDriverOptions {
	values := map[string]interface{}{
		"nutanix-username":   "user",
		"nutanix-password":   "pass",
		"nutanix-endpoint":   "pc.example.com",
		"nutanix-cluster":    "cluster1",
		"nutanix-vm-network": []string{"vlan0"},
		"nutanix-vm-image":   "image1",
	}
	for k, v := range overrides {
		values[k] = v
	}

	return &drivers.CheckDriverOptions{
		FlagsValues: values,
		CreateFlags: NewDriver("default", "path").GetCreateFlags(),
	}
}

func TestSetConfigFromFlags_LinuxDefault(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(nil)

	if err := d.SetConfigFromFlags(checkFlags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(checkFlags.InvalidFlags) != 0 {
		t.Errorf("unexpected invalid flags: %v", checkFlags.InvalidFlags)
	}
	if d.OS != "linux" {
		t.Errorf("expected default OS linux, got %q", d.OS)
	}
	if d.SSHUser != "root" {
		t.Errorf("expected default linux ssh user root, got %q", d.SSHUser)
	}
}

func TestSetConfigFromFlags_WindowsDefaultsSSHUserToAdministrator(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-os": "windows",
	})

	if err := d.SetConfigFromFlags(checkFlags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(checkFlags.InvalidFlags) != 0 {
		t.Errorf("unexpected invalid flags: %v", checkFlags.InvalidFlags)
	}
	if d.OS != "windows" {
		t.Errorf("expected OS windows, got %q", d.OS)
	}
	if d.SSHUser != "Administrator" {
		t.Errorf("expected default windows ssh user Administrator, got %q", d.SSHUser)
	}
}

func TestSetConfigFromFlags_SSHUserOverrideWins(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-os":       "windows",
		"nutanix-vm-ssh-user": "deployer",
	})

	if err := d.SetConfigFromFlags(checkFlags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.SSHUser != "deployer" {
		t.Errorf("expected explicit ssh user override deployer, got %q", d.SSHUser)
	}
}

func TestSetConfigFromFlags_InvalidOS(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-os": "bogus",
	})

	err := d.SetConfigFromFlags(checkFlags)
	if err == nil {
		t.Fatal("expected an error for an invalid nutanix-vm-os value, got nil")
	}
}

func TestSetConfigFromFlags_StaticIPUnsetByDefault(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(nil)

	if err := d.SetConfigFromFlags(checkFlags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.VMIP != "" {
		t.Errorf("expected VMIP to default empty, got %q", d.VMIP)
	}
}

func TestSetConfigFromFlags_ValidStaticIP(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-ip": "10.0.0.50",
	})

	if err := d.SetConfigFromFlags(checkFlags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(checkFlags.InvalidFlags) != 0 {
		t.Errorf("unexpected invalid flags: %v", checkFlags.InvalidFlags)
	}
	if d.VMIP != "10.0.0.50" {
		t.Errorf("expected VMIP 10.0.0.50, got %q", d.VMIP)
	}
}

func TestSetConfigFromFlags_InvalidStaticIP(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-ip": "not-an-ip",
	})

	if err := d.SetConfigFromFlags(checkFlags); err == nil {
		t.Fatal("expected an error for an invalid nutanix-vm-ip value, got nil")
	}
}

func TestSetConfigFromFlags_IPv6StaticIPRejected(t *testing.T) {
	d := NewDriver("default", "path")
	checkFlags := baseValidFlags(map[string]interface{}{
		"nutanix-vm-ip": "2001:db8::1",
	})

	if err := d.SetConfigFromFlags(checkFlags); err == nil {
		t.Fatal("expected an error for an IPv6 nutanix-vm-ip value (API only supports IPv4), got nil")
	}
}
