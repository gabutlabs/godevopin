<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <div>
    <v-toolbar
      density="comfortable"
      :color="appStore.theme == 'light' ? 'white' : '#2C303A'"
      title="Project Detail"
    >
      <template #prepend>
        <v-btn icon="mdi-arrow-left" @click="router.back()"></v-btn>
      </template>
    </v-toolbar>

    <v-container fluid>
      <v-row v-if="projectStore.project">
        <v-col cols="12">
          <v-card variant="outlined">
            <v-card-text>
              <v-row>
                <v-col cols="12" md="3">
                  <div class="text-caption text-medium-emphasis">
                    Project Name
                  </div>
                  <div class="text-h6">{{ projectStore.project.name }}</div>
                </v-col>
                <v-col cols="12" md="3">
                  <div class="text-caption text-medium-emphasis">Type</div>
                  <v-chip size="small" color="primary" label>{{
                    projectStore.project.project_type
                  }}</v-chip>
                </v-col>
                <v-col cols="12" md="6">
                  <div class="text-caption text-medium-emphasis">Log Path</div>
                  <div class="text-body-2 font-weight-medium">
                    {{ projectStore.project.path_log }}
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12">
          <v-card variant="outlined">
            <v-card-title class="d-flex align-center py-4 px-6">
              <v-icon start icon="mdi-history"></v-icon>
              Log History
              <v-spacer></v-spacer>
              <v-btn
                icon="mdi-refresh"
                variant="text"
                @click="loadLogs"
                :loading="projectStore.loading"
              ></v-btn>
            </v-card-title>
            <v-divider></v-divider>
            <v-data-table-server
              v-model:items-per-page="itemsPerPage"
              :headers="headers"
              :items="projectStore.logs"
              :items-length="projectStore.logsTotal"
              :loading="projectStore.loading"
              item-value="id"
              @update:options="loadLogs"
              density="compact"
            >
              <template #[`item.timestamp`]="{ item }">
                <span class="text-caption">{{
                  formatDate(item.timestamp)
                }}</span>
              </template>
              <template #[`item.level`]="{ item }">
                <v-chip
                  :color="getLevelColor(item.level)"
                  size="x-small"
                  label
                  class="font-weight-bold text-uppercase"
                >
                  {{ item.level }}
                </v-chip>
              </template>
              <template #[`item.message`]="{ item }">
                <div
                  class="text-truncate"
                  style="max-width: 600px"
                  :title="item.message"
                >
                  {{ item.message }}
                </div>
              </template>
              <template #[`item.source`]="{ item }">
                <span class="text-caption text-medium-emphasis">{{
                  item.source
                }}</span>
              </template>
            </v-data-table-server>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useProjectStore } from "@/stores/project";
import { useAppStore } from "@/stores/app";

const route = useRoute();
const router = useRouter();
const projectStore = useProjectStore();
const appStore = useAppStore();

const id = Number((route.params as any).id);
const itemsPerPage = ref(50);

const headers = [
  { title: "Timestamp", key: "timestamp", width: "200px" },
  { title: "Level", key: "level", width: "100px" },
  { title: "Message", key: "message" },
  { title: "Source", key: "source", width: "250px" },
];

const loadLogs = async (options: any = { page: 1, itemsPerPage: 50 }) => {
  await projectStore.fetchProjectLogs(id, options.page, options.itemsPerPage);
};

const getLevelColor = (level: string) => {
  if (!level) return "grey";
  const l = level.toUpperCase();
  if (l.includes("ERR") || l.includes("FAIL") || l.includes("CRIT"))
    return "error";
  if (l.includes("WARN")) return "warning";
  if (l.includes("INFO")) return "info";
  if (l.includes("DEB")) return "secondary";
  return "grey";
};

const formatDate = (date: string) => {
  if (!date) return "-";
  return new Date(date).toLocaleString();
};

onMounted(async () => {
  await projectStore.fetchProject(id);
  await loadLogs();
});
</script>
