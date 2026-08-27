package driver

import (
	"strings"
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	v3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/nutanix-cloud-native/prism-go-client/v3/models"
	"github.com/nutanix/docker-machine/utils"
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

func TestIsUUID(t *testing.T) {
	cases := map[string]bool{
		"550e8400-e29b-41d4-a716-446655440000": true,
		"550E8400-E29B-41D4-A716-446655440000": true,
		"my-cluster":                           false,
		"":                                     false,
		"550e8400-e29b-41d4-a716":              false,
	}

	for input, want := range cases {
		if got := isUUID(input); got != want {
			t.Errorf("isUUID(%q) = %v, want %v", input, got, want)
		}
	}
}

func clusterEntity(name, uuid string) *v3.ClusterIntentResponse {
	return &v3.ClusterIntentResponse{
		Spec:     &models.Cluster{Name: name},
		Metadata: &v3.Metadata{UUID: utils.StringPtr(uuid)},
	}
}

func TestFindClusterUUID_SingleMatch(t *testing.T) {
	entities := []*v3.ClusterIntentResponse{
		clusterEntity("other", "11111111-1111-1111-1111-111111111111"),
		clusterEntity("mycluster", "22222222-2222-2222-2222-222222222222"),
	}

	got, err := findClusterUUID("mycluster", entities)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("expected the matching cluster's UUID, got %q", got)
	}
}

func TestFindClusterUUID_NoMatch(t *testing.T) {
	entities := []*v3.ClusterIntentResponse{
		clusterEntity("other", "11111111-1111-1111-1111-111111111111"),
	}

	_, err := findClusterUUID("mycluster", entities)
	if err == nil {
		t.Fatal("expected an error when no cluster matches the name, got nil")
	}
}

// Nutanix does not enforce cluster name uniqueness within a Prism Central,
// so this is a real, user-hit scenario (see rancher/rancher#47493) - the
// error must name both conflicting UUIDs so nutanix-cluster can be set to
// one of them to disambiguate, instead of failing the same way forever.
func TestFindClusterUUID_MultipleMatchesListsUUIDsForDisambiguation(t *testing.T) {
	entities := []*v3.ClusterIntentResponse{
		clusterEntity("mycluster", "11111111-1111-1111-1111-111111111111"),
		clusterEntity("mycluster", "22222222-2222-2222-2222-222222222222"),
	}

	_, err := findClusterUUID("mycluster", entities)
	if err == nil {
		t.Fatal("expected an error when more than one cluster matches the name, got nil")
	}
	for _, uuid := range []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"} {
		if !strings.Contains(err.Error(), uuid) {
			t.Errorf("expected error to mention conflicting UUID %s so the user can disambiguate, got: %v", uuid, err)
		}
	}
}
