<template>
  <page-content title="Worker Services">
    <template #append>
      <v-btn prepend-icon="mdi-plus" variant="outlined" @click="dialog = true"
        >Add Worker Service</v-btn
      >
    </template>
    <v-card>
      <v-card-text>
        <v-text-field
          v-model="name"
          class="ma-2"
          density="compact"
          placeholder="Search worker service..."
          width="30%"
        ></v-text-field>
        <v-data-table :items="state.workerServices" :headers="headers">
          <template #item.actions="{ item }">
            <ActionTable
              :action-items="getActionMenuItems(item)"
              :item="item"
            />
          </template>
          <template #item.desired_state="{ item }">
            <v-chip
              :color="item.desired_state === 'enabled' ? 'success' : 'warning'"
              size="small"
              variant="elevated"
            >
              {{ item.desired_state }}
            </v-chip>
          </template>
          <template #item.current_status="{ item }">
            <v-chip
              :color="getStatusColor(item.current_status)"
              size="small"
              variant="elevated"
            >
              {{ item.current_status }}
            </v-chip>
          </template>
          <template #item.health_status="{ item }">
            <v-chip
              :color="getHealthColor(item.health_status)"
              size="small"
              variant="elevated"
            >
              {{ item.health_status }}
            </v-chip>
          </template>
          <template #item.last_success_at="{ item }">
            {{
              item.last_success_at
                ? new Date(item.last_success_at).toLocaleString()
                : ""
            }}
          </template>
          <template #item.created_at="{ item }">
            {{ new Date(item.created_at).toLocaleString() }}
          </template>
          <template #item.updated_at="{ item }">
            {{ new Date(item.updated_at).toLocaleString() }}
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>

    <!-- Worker Service Form Dialog -->
    <v-dialog v-model="dialog" width="auto" max-width="700">
      <v-card
        max-width="700"
        width="700"
        prepend-icon="mdi-cog-transfer"
        title="Form Worker Service"
      >
        <v-card-text>
          <v-form @submit.prevent="submit">
            <v-row>
              <v-col xl="6" md="6" sm="12">
                <v-text-field
                  v-model="payload.name"
                  label="Name"
                  :error-messages="v$.name.$errors.map((e: { $message: any; }) => e.$message)"
                  @blur="v$.name.$touch"
                  @input="v$.name.$touch"
                ></v-text-field>
              </v-col>
              <v-col xl="6" md="6" sm="12">
                <v-select
                  v-model="payload.desired_state"
                  :items="['enabled', 'disabled']"
                  label="Desired State"
                  :error-messages="v$.desired_state.$errors.map((e: { $message: any; }) => e.$message)"
                  @blur="v$.desired_state.$touch"
                  @input="v$.desired_state.$touch"
                ></v-select>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="payload.description"
                  label="Description"
                  :error-messages="v$.description.$errors.map((e: { $message: any; }) => e.$message)"
                  @blur="v$.description.$touch"
                  @input="v$.description.$touch"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-btn class="mt-2" type="submit" block variant="elevated"
              >Submit</v-btn
            >
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </page-content>
</template>

<script lang="ts" setup>
import { useNotify } from "@/composables/useNotify";
import {
  useWorkerServiceStore,
  type WorkerService,
  type CreateWorkerServiceRequest,
  type UpdateWorkerServiceRequest,
} from "@/stores/worker-service";
import useVuelidate from "@vuelidate/core";
import { required, maxLength } from "@vuelidate/validators";
import { computed, ref, watch, reactive, onMounted } from "vue";
const notify = useNotify();
const state = useWorkerServiceStore();
const dialog = ref(false);
const name = ref("");

// Initialize payload with reactive data
const payload = reactive({
  id: 0,
  name: "",
  description: "",
  desired_state: "enabled" as "enabled" | "disabled",
  is_updated: false,
} as CreateWorkerServiceRequest & { id: number; is_updated: boolean });

// Define validation rules
const rules = {
  name: { required, maxLength: maxLength(255) },
  description: { maxLength: maxLength(1000) },
  desired_state: { required },
};

const v$ = useVuelidate(rules, payload);

// Load worker services on component mount
onMounted(async () => {
  await state.fetchAllWorkerServices();
});

// Watch search input to filter worker services
watch(name, async (newName) => {
  await state.setFilterWorkerService(newName);
});

// Define table headers
const headers = computed(() => [
  { title: "Name", key: "name" },
  { title: "Description", key: "description" },
  { title: "Desired State", key: "desired_state" },
  { title: "Current Status", key: "current_status" },
  { title: "Health Status", key: "health_status" },
  { title: "PID", key: "pid" },
  { title: "Last Heartbeat", key: "last_heartbeat_at" },
  { title: "Last Success", key: "last_success_at" },
  { title: "Created At", key: "created_at" },
  { title: "Updated At", key: "updated_at" },
  { title: "Actions", key: "actions", sortable: false },
]);

// Define method to get action menu items with proper icons
const getActionMenuItems = (item: WorkerService) => [
  {
    title: "Edit",
    icon: "mdi-pencil",
    onClick: (item: WorkerService) => {
      payload.id = item.id;
      payload.name = item.name;
      payload.description = item.description;
      payload.is_updated = true;
      dialog.value = true;
    },
  },
  {
    title: "Start",
    icon: "mdi-play",
    disabled: (item: WorkerService) =>
      item.current_status === "running" || item.current_status === "starting",
    onClick: async (item: WorkerService) => {
      await state.updateStatusWorkerService(item.id, "starting");
      if (state.action_result.is_success) {
        await state.fetchAllWorkerServices();
        notify.success(`Worker service started successfully`);
      } else {
        notify.error(state.action_result.message);
      }
    },
  },
  {
    title: "Stop",
    icon: "mdi-stop",
    disabled: (item: WorkerService) => item.current_status === "stopped",
    onClick: async (item: WorkerService) => {
      await state.updateStatusWorkerService(item.id, "stopped");
      if (state.action_result.is_success) {
        await state.fetchAllWorkerServices();
        notify.success(`Worker service stopped successfully`);
      } else {
        notify.error(state.action_result.message);
      }
    },
  },
  {
    title: "Restart",
    icon: "mdi-restart",
    disabled: (item: WorkerService) => item.current_status === "starting",
    onClick: async (item: WorkerService) => {
      await state.updateStatusWorkerService(item.id, "restart");
      if (state.action_result.is_success) {
        await state.fetchAllWorkerServices();
        notify.success(`Worker service restarted successfully`);
      } else {
        notify.error(state.action_result.message);
      }
    },
  },
  {
    title: "Delete",
    style: "color: red",
    icon: "mdi-delete",
    onClick: async (item: WorkerService) => {
      await state.deleteWorkerService(item.id);
      if (state.action_result.is_success) {
        await state.fetchAllWorkerServices();
        notify.success(state.action_result.message);
      } else {
        notify.error(state.action_result.message);
      }
    },
  },
];

// Helper function to get color for status
const getStatusColor = (status: string) => {
  switch (status) {
    case "running":
      return "success";
    case "stopped":
      return "secondary";
    case "failed":
      return "error";
    case "degraded":
      return "warning";
    case "starting":
      return "info";
    default:
      return "default";
  }
};

// Helper function to get color for health status
const getHealthColor = (status: string) => {
  switch (status) {
    case "healthy":
      return "success";
    case "unhealthy":
      return "error";
    default:
      return "secondary";
  }
};

// Form submission handler
async function submit() {
  if (!(await v$.value.$validate())) {
    return;
  }

  if (payload.is_updated) {
    await state.updateWorkerService(
      payload.id,
      payload as UpdateWorkerServiceRequest
    );
  } else {
    await state.createWorkerService(payload);
  }

  if (state.action_result.is_success) {
    notify.success(state.action_result.message);
    dialog.value = false;

    // Reset form after successful submission
    resetForm();

    await state.fetchAllWorkerServices();
  } else {
    notify.error(
      `${state.action_result.message} : \n ${
        state.action_result.data != null
          ? Object.values(state.action_result.data).join("\n ")
          : ""
      }`
    );
  }
}

// Reset form function
function resetForm() {
  payload.id = 0;
  payload.name = "";
  payload.description = "";
  payload.desired_state = "enabled";
  payload.is_updated = false;
  v$.value.$reset();
}
</script>
