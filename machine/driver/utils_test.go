package driver

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBuildBoolNodes(t *testing.T) {
	cases := []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, c := range cases {
		nodes := buildBoolNodes("inactive", c.value)
		if len(nodes) != 2 {
			t.Fatalf("expected 2 nodes, got %d", len(nodes))
		}
		if nodes[0].Value != "inactive" {
			t.Errorf("expected key node value %q, got %q", "inactive", nodes[0].Value)
		}
		if nodes[1].Tag != "!!bool" {
			t.Errorf("expected value node tag !!bool, got %q", nodes[1].Tag)
		}
		if nodes[1].Value != c.expected {
			t.Errorf("expected value node value %q, got %q", c.expected, nodes[1].Value)
		}
	}
}

func TestBuildCloudInitUser_Linux(t *testing.T) {
	node := buildCloudInitUser("root", []byte("ssh-rsa AAAAtest"), false)

	out, err := yaml.Marshal(node)
	if err != nil {
		t.Fatalf("failed to marshal node: %v", err)
	}
	rendered := string(out)

	if !strings.Contains(rendered, "name: root") {
		t.Errorf("expected rendered user to contain name: root, got: %s", rendered)
	}
	if !strings.Contains(rendered, "sudo: ALL=(ALL) NOPASSWD:ALL") {
		t.Errorf("expected rendered user to contain sudo entry, got: %s", rendered)
	}
	if strings.Contains(rendered, "inactive") {
		t.Errorf("did not expect linux user to contain an inactive entry, got: %s", rendered)
	}
	if !strings.Contains(rendered, "ssh-rsa AAAAtest") {
		t.Errorf("expected rendered user to contain the ssh public key, got: %s", rendered)
	}
}

func TestBuildCloudInitUser_Windows(t *testing.T) {
	node := buildCloudInitUser("Administrator", []byte("ssh-rsa AAAAtest"), true)

	out, err := yaml.Marshal(node)
	if err != nil {
		t.Fatalf("failed to marshal node: %v", err)
	}
	rendered := string(out)

	if !strings.Contains(rendered, "name: Administrator") {
		t.Errorf("expected rendered user to contain name: Administrator, got: %s", rendered)
	}
	if !strings.Contains(rendered, "inactive: false") {
		t.Errorf("expected rendered user to contain inactive: false, got: %s", rendered)
	}
	if strings.Contains(rendered, "sudo") {
		t.Errorf("did not expect windows user to contain a sudo entry, got: %s", rendered)
	}
	if !strings.Contains(rendered, "ssh-rsa AAAAtest") {
		t.Errorf("expected rendered user to contain the ssh public key, got: %s", rendered)
	}
}

func TestDefaultCloudInitUserData_LinuxMatchesLegacyDefault(t *testing.T) {
	pubKey := []byte("ssh-rsa AAAAtest")

	got := defaultCloudInitUserData("root", pubKey, false)
	want := []byte("#cloud-config\r\nusers:\r\n - name: root\r\n   ssh_authorized_keys:\r\n    - " + string(pubKey))

	if string(got) != string(want) {
		t.Errorf("linux default cloud-init changed from the pre-existing hardcoded default.\ngot:  %q\nwant: %q", got, want)
	}
}

func TestDefaultCloudInitUserData_Windows(t *testing.T) {
	pubKey := []byte("ssh-rsa AAAAtest")

	got := string(defaultCloudInitUserData("Administrator", pubKey, true))

	if !strings.Contains(got, "name: Administrator") {
		t.Errorf("expected default windows cloud-init to contain name: Administrator, got: %s", got)
	}
	if !strings.Contains(got, "inactive: false") {
		t.Errorf("expected default windows cloud-init to contain inactive: false, got: %s", got)
	}
	if !strings.HasPrefix(got, "#cloud-config") {
		t.Errorf("expected default windows cloud-init to start with #cloud-config, got: %s", got)
	}
}
