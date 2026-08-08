<script>
import LabeledInput from '@shell/components/form/LabeledInput';
import Checkbox from '@shell/components/form/Checkbox';

// Field names below (endpoint/port/username/password/insecure) mirror the
// `nutanix-*` driver flags in machine/driver/driver.go with the `nutanix-`
// prefix stripped and the remainder camelCased. Verify this mapping against
// a live `yarn dev` session before shipping - see the package README.
export default {
  components: { LabeledInput, Checkbox },

  props: {
    value: {
      type:    Object,
      default: () => {},
    },

    mode: {
      type:    String,
      default: 'create',
    },
  },

  data() {
    if (!this.value.decodedData) {
      this.$set(this.value, 'decodedData', {});
    }
    if (this.value.decodedData.port === undefined) {
      this.$set(this.value.decodedData, 'port', '9440');
    }

    return {};
  },

  methods: {
    // Rancher's cloud credential form calls into a `test`-style hook on
    // some driver integrations to gate the "Create" button; wire this up
    // per whatever convention your target Rancher version documents, this
    // is a placeholder for basic required-field validation.
    validate() {
      const data = this.value.decodedData || {};

      return !!(data.endpoint && data.username && data.password);
    },
  },
};
</script>

<template>
  <div>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.decodedData.endpoint"
          label="Prism Central Endpoint"
          placeholder="pc.example.com"
          required
        />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model="value.decodedData.port"
          label="Port"
          placeholder="9440"
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-6">
        <LabeledInput
          v-model="value.decodedData.username"
          label="Username"
          placeholder="admin@example.com, or X-ntnx-api-key for a service account"
          required
        />
      </div>
      <div class="col span-6">
        <LabeledInput
          v-model="value.decodedData.password"
          label="Password / API Key"
          type="password"
          required
        />
      </div>
    </div>
    <div class="row mb-10">
      <div class="col span-6">
        <Checkbox
          v-model="value.decodedData.insecure"
          label="Allow insecure SSL connections"
        />
      </div>
    </div>
  </div>
</template>
