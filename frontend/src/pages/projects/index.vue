<style scoped lang="scss"></style>
<template>
  <page-content title="Projects">
    <template #append>
      <v-btn prepend-icon="mdi-plus" variant="flat" color="primary" @click="openAddDialog"
        >Add Project</v-btn
      >
    </template>
    <v-card>
      <v-card-text>
        <v-text-field
          v-model="search"
          class="ma-2"
          density="compact"
          placeholder="Search project..."
          variant="outlined"
          width="30%"
        ></v-text-field>
        <v-data-table :items="filteredProjects" :headers="headers" :loading="state.loading">
          <template #item.actions="{ item }">
            <ActionTable :action-items="actionMenuItems" :item="item" />
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
    <v-dialog v-model="dialog" width="auto" max-width="700">
      <v-card
        max-width="700"
        width="700"
        prepend-icon="mdi-folder-outline"
        title="Form Project"
      >
        <v-card-text>
          <v-form @submit.prevent="submit">
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="payload.name"
                  label="Project Name"
                  :error-messages="v$.name.$errors.map((e: any) => e.$message)"
                  @blur="v$.name.$touch"
                  @input="v$.name.$touch"
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="payload.path_log"
                  label="Log Path (e.g. /home/user/project)"
                  :error-messages="v$.path_log.$errors.map((e: any) => e.$message)"
                  @blur="v$.path_log.$touch"
                  @input="v$.path_log.$touch"
                  hint="Absolute path to the log file or directory"
                  persistent-hint
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="payload.project_type"
                  label="Project Type"
                  :items="['laravel', 'nodejs', 'python', 'custom']"
                  :error-messages="v$.project_type.$errors.map((e: any) => e.$message)"
                  @blur="v$.project_type.$touch"
                  @input="v$.project_type.$touch"
                ></v-select>
              </v-col>
              <v-col cols="12" v-if="payload.project_type !== 'laravel'">
                <div class="d-flex align-center justify-space-between mb-2">
                  <span class="text-subtitle-1">Log Format Builder</span>
                  <v-btn size="small" color="primary" variant="outlined" prepend-icon="mdi-plus" @click="addComponent">Add Component</v-btn>
                </div>
                <v-row v-for="(comp, index) in log_components" :key="index" dense align="center" class="mb-1">
                  <v-col cols="5">
                    <v-select
                      v-model="comp.type"
                      :items="availableComponents"
                      label="Type"
                      density="compact"
                      hide-details
                    ></v-select>
                  </v-col>
                  <v-col cols="5" v-if="comp.type === 'separator'">
                    <v-combobox
                      v-model="comp.value"
                      :items="availableSeparators"
                      label="Separator/Text"
                      density="compact"
                      hide-details
                    ></v-combobox>
                  </v-col>
                  <v-col cols="2">
                    <v-btn icon="mdi-delete" size="small" color="error" variant="text" @click="removeComponent(index)"></v-btn>
                  </v-col>
                </v-row>
                <v-alert type="info" variant="tonal" class="mt-2" density="compact" v-if="log_components.length > 0">
                  <div class="text-caption">Preview Pattern:</div>
                  <code>{{ formatPreview }}</code>
                </v-alert>
              </v-col>
            </v-row>
            <v-btn class="mt-4" type="submit" block variant="flat" color="primary"
              >Submit</v-btn
            >
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </page-content>
</template>
<script setup lang="ts">
import { useNotify } from "@/composables/useNotify";
import { useProjectStore, type Project } from "@/stores/project";
import useVuelidate from "@vuelidate/core";
import { required } from "@vuelidate/validators";
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useRouter } from "vue-router";

const search = ref("");
const dialog = ref(false);
const state = useProjectStore();
const router = useRouter();
const payload = reactive({
  id: 0,
  name: "",
  path_log: "",
  project_type: "laravel",
  log_format: "",
  is_updated: false,
});

const log_components = ref<Array<{ type: string, value: string }>>([]);
const availableComponents = [
  { title: 'Timestamp', value: 'timestamp' },
  { title: 'Log Level', value: 'level' },
  { title: 'Message', value: 'message' },
  { title: 'Source/File', value: 'source' },
  { title: 'Separator/Text', value: 'separator' },
];
const availableSeparators = [' ', ' - ', ' | ', ' : ', ' [ ', ' ] ', ' ( ', ' ) '];

function addComponent() {
  log_components.value.push({ type: 'timestamp', value: '' });
}

function removeComponent(index: number) {
  log_components.value.splice(index, 1);
}

const formatPreview = computed(() => {
  return log_components.value.map(c => {
    if (c.type === 'separator') return c.value;
    return `{${c.type}}`;
  }).join('');
});

watch(formatPreview, (newVal) => {
  payload.log_format = newVal;
});

function parseLogFormat(format: string) {
  if (!format) return [];
  const components: Array<{ type: string, value: string }> = [];
  const parts = format.split(/(\{.*?\})/);
  parts.forEach(part => {
    if (!part) return;
    if (part.startsWith('{') && part.endsWith('}')) {
      const type = part.substring(1, part.length - 1);
      components.push({ type, value: '' });
    } else {
      components.push({ type: 'separator', value: part });
    }
  });
  return components;
}

const notify = useNotify();
const rules = {
  name: { required },
  path_log: { required },
  project_type: { required },
};
const v$ = useVuelidate(rules, payload);

const filteredProjects = computed(() => {
  if (!search.value) return state.projects;
  return state.projects.filter(p => p.name.toLowerCase().includes(search.value.toLowerCase()));
});

function openAddDialog() {
  payload.id = 0;
  payload.name = "";
  payload.path_log = "";
  payload.project_type = "laravel";
  payload.log_format = "";
  payload.is_updated = false;
  log_components.value = [];
  v$.value.$reset();
  dialog.value = true;
}

async function submit() {
  if (!(await v$.value.$validate())) {
    return;
  }
  
  const submitPayload = {
    name: payload.name,
    path_log: payload.path_log,
    project_type: payload.project_type,
    log_format: payload.log_format
  };

  if (payload.is_updated) {
    await state.updateProject(payload.id, submitPayload);
  } else {
    await state.createProject(submitPayload);
  }
  
  if (state.action_result.is_success) {
    notify.success(state.action_result.message);
    dialog.value = false;
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

const actionMenuItems = [
  {
    title: "View Details",
    icon: "mdi-eye",
    onClick: (item: Project) => {
      router.push(`/projects/${item.id}`);
    },
  },
  {
    title: "Edit",
    icon: "mdi-pencil",
    onClick: (item: Project) => {
      payload.id = item.id;
      payload.name = item.name;
      payload.path_log = item.path_log;
      payload.project_type = item.project_type;
      payload.log_format = item.log_format;
      payload.is_updated = true;
      log_components.value = parseLogFormat(item.log_format);
      v$.value.$reset();
      dialog.value = true;
    },
  },
  {
    title: "Delete",
    style: "color: red",
    icon: "mdi-delete",
    onClick: async (item: Project) => {
      if (confirm(`Are you sure you want to delete project ${item.name}?`)) {
        await state.deleteProject(item.id);
        if (state.action_result.is_success) {
          notify.success(state.action_result.message);
        } else {
          notify.error(state.action_result.message);
        }
      }
    },
  },
];

const headers = [
  { title: "Name", value: "name" },
  { title: "Log Path", value: "path_log" },
  { title: "Type", value: "project_type" },
  { title: "Created At", value: "created_at" },
  { title: "Actions", value: "actions", sortable: false },
];

onMounted(async () => {
  await state.fetchAllProjects();
});
</script>
