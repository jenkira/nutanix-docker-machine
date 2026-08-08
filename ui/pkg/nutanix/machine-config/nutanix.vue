<script>
import LabeledInput from '@shell/components/form/LabeledInput';
import LabeledSelect from '@shell/components/form/LabeledSelect';
import Checkbox from '@shell/components/form/Checkbox';

const OS_OPTIONS = [
  { label: 'Linux', value: 'linux' },
  { label: 'Windows', value: 'windows' },
];

const BOOT_TYPE_OPTIONS = [
  { label: 'Legacy (BIOS)', value: 'legacy' },
  { label: 'UEFI', value: 'uefi' },
];

// Field names below mirror the `nutanix-*` driver flags in
// machine/driver/driver.go with the `nutanix-` prefix stripped and the
// remainder camelCased (e.g. `nutanix-vm-os` -> `vmOs`). Verify this mapping
// against a live `yarn dev` session before shipping - see the package README.
export default {
  components: { LabeledInput, LabeledSelect, Checkbox },

  props: {
    // The machine config resource; fields are bound directly onto it.
    value: {
      type:    Object,
      default: () => {},
    },

    // Id of the cloud credential selected earlier in the wizard.
    credentialId: {
      type:    String,
      default: '',
    },

    mode: {
      type:    String,
      default: 'create',
    },
  },

  data() {
    // Mirror the driver's own defaults (see defaultOSType/defaultBootType/
    // defaultVMMem/defaultVCPUs/defaultCores in driver.go) so the form
    // reflects reality before the user touches anything.
    const defaults = {
      vmOs: 'linux', bootType: 'legacy', vmMem: 2048, vmCpus: 2, vmCores: 1,
    };

    Object.keys(defaults).forEach((key) => {
      if (this.value[key] === undefined) {
        this.$set(this.value, key, defaults[key]);
      }
    });

    return {
      osOptions:       OS_OPTIONS,
      bootTypeOptions: BOOT_TYPE_OPTIONS,
    };
  },

  computed: {
    isWindows() {
      return this.value.vmOs === 'windows';
    },

    sshUserPlaceholder() {
      return this.isWindows ? 'Administrator' : 'root';
    },
  },
};
</script>

<template>
  <div>
    <h3>Target</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.cluster"
          label="Cluster"
          placeholder="Cluster name (case sensitive)"
          required
        />
      </div>
      <div class="col span-6">
        <LabeledSelect
          v-model="value.vmNetwork"
          label="Network(s)"
          :taggable="true"
          :multiple="true"
          :options="[]"
          placeholder="Network name or UUID - press enter to add"
        />
      </div>
    </div>

    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.vmImage"
          label="Image"
          placeholder="Disk image template name"
          required
        />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model.number="value.vmImageSize"
          label="Image size override (GiB)"
          type="number"
        />
      </div>
    </div>

    <h3>Operating System</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledSelect
          v-model="value.vmOs"
          label="Guest OS"
          :options="osOptions"
        />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model="value.vmSshUser"
          label="SSH User"
          :placeholder="sshUserPlaceholder"
        />
      </div>
    </div>
    <p v-if="isWindows" class="text-muted mb-10">
      The image must already have cloudbase-init, an SSH server, and Nutanix VirtIO
      drivers installed - see the driver README's "Windows Node Support" section.
    </p>

    <h3>Sizing</h3>
    <div class="row mb-10">
      <div class="col span-4">
        <LabeledInput
          v-model.number="value.vmMem"
          label="Memory (MB)"
          type="number"
        />
      </div>
      <div class="col span-4">
        <LabeledInput
          v-model.number="value.vmCpus"
          label="vCPUs"
          type="number"
        />
      </div>
      <div class="col span-4">
        <LabeledInput
          v-model.number="value.vmCores"
          label="Cores per vCPU"
          type="number"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-6">
        <Checkbox
          v-model="value.vmCpuPassthrough"
          label="Passthrough host CPU features"
        />
      </div>
      <div class="col span-6">
        <LabeledSelect
          v-model="value.bootType"
          label="Boot Type"
          :options="bootTypeOptions"
        />
      </div>
    </div>

    <h3>Storage</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.storageContainer"
          label="Additional disk storage container UUID"
        />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model.number="value.diskSize"
          label="Additional disk size (GiB)"
          type="number"
        />
      </div>
    </div>

    <h3>Placement &amp; Metadata</h3>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.project"
          label="Project"
          placeholder="default"
        />
      </div>
      <div class="col span-6">
        <LabeledSelect
          v-model="value.vmCategories"
          label="Categories"
          :taggable="true"
          :multiple="true"
          :options="[]"
          placeholder="key=value - press enter to add"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledSelect
          v-model="value.vmGpu"
          label="GPU device names"
          :taggable="true"
          :multiple="true"
          :options="[]"
        />
      </div>
      <div class="col span-6">
        <Checkbox
          v-model="value.vmSerialPort"
          label="Attach serial port"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-12">
        <LabeledInput
          v-model="value.vmDescription"
          label="VM Description"
          placeholder="VM created by Nutanix Rancher Node Driver"
        />
      </div>
    </div>

    <h3>Cloud-Init</h3>
    <div class="row mb-10">
      <div class="col span-12">
        <LabeledInput
          v-model="value.cloudInit"
          label="Cloud-init / cloudbase-init user-data"
          type="multiline"
          :min-height="160"
          placeholder="#cloud-config..."
        />
      </div>
    </div>
  </div>
</template>
