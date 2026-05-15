<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <div>
    <PageContent title="Application Settings">
    </PageContent>

    <v-container class="mt-4">
      <v-row>
        <v-col cols="12" md="8" lg="6">
          <v-card :loading="settingStore.loading">
            <v-card-title class="d-flex align-center">
              <v-icon start icon="mdi-cog"></v-icon>
              General Configuration
            </v-card-title>
            <v-divider></v-divider>
            <v-card-text>
              <v-form @submit.prevent="saveSettings">
                <v-row>
                  <v-col cols="12">
                    <div class="text-subtitle-1 mb-2">Monitoring Intervals</div>
                    <v-text-field
                      v-model.number="form.monitoring_interval_seconds"
                      label="Monitoring Interval (seconds)"
                      type="number"
                      hint="Interval for system metrics collection"
                      persistent-hint
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model.number="form.alarm_check_interval_seconds"
                      label="Alarm Check Interval (seconds)"
                      type="number"
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model.number="form.alarm_repeat_interval_minutes"
                      label="Alarm Repeat Interval (minutes)"
                      type="number"
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12">
                    <v-divider class="my-4"></v-divider>
                    <div class="text-subtitle-1 mb-2">Critical Thresholds</div>
                  </v-col>

                  <v-col cols="12" sm="4">
                    <v-text-field
                      v-model.number="form.system_cpu_critical_percent"
                      label="CPU Critical (%)"
                      type="number"
                      suffix="%"
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12" sm="4">
                    <v-text-field
                      v-model.number="form.system_disk_critical_percent"
                      label="Disk Critical (%)"
                      type="number"
                      suffix="%"
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12" sm="4">
                    <v-text-field
                      v-model.number="form.system_mem_critical_percent"
                      label="Memory Critical (%)"
                      type="number"
                      suffix="%"
                      variant="outlined"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="12">
                    <v-divider class="my-4"></v-divider>
                    <div class="text-subtitle-1 mb-2">Worker Settings</div>
                  </v-col>

                  <v-col cols="12">
                    <v-text-field
                      v-model.number="form.worker_heartbeat_timeout_seconds"
                      label="Worker Heartbeat Timeout (seconds)"
                      type="number"
                      variant="outlined"
                      hint="Timeout before a worker is considered inactive"
                      persistent-hint
                    ></v-text-field>
                  </v-col>
                </v-row>
              </v-form>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                variant="flat"
                @click="saveSettings"
                :loading="settingStore.loading"
                prepend-icon="mdi-content-save"
              >
                Save Settings
              </v-btn>
            </v-card-actions>
          </v-card>

          <v-alert
            v-if="settingStore.actionResult.message"
            :type="settingStore.actionResult.isSuccess ? 'success' : 'error'"
            class="mt-4"
            closable
          >
            {{ settingStore.actionResult.message }}
          </v-alert>
        </v-col>

        <v-col cols="12" md="4">
          <v-card>
            <v-card-title>Information</v-card-title>
            <v-card-text>
              <p class="text-body-2 mb-4">
                These settings control the behavior of the monitoring system and background workers.
              </p>
              <p class="text-caption text-medium-emphasis">
                Last updated: {{ settingStore.settings?.updated_at ? new Date(settingStore.settings.updated_at).toLocaleString() : 'Never' }}
              </p>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts" setup>
import { reactive, onMounted, watch } from 'vue';
import { useAppStore } from "@/stores/app";
import { useSettingStore, type UpdateSettingRequest } from "@/stores/setting";
import PageContent from "@/components/PageContent.vue";

const appStore = useAppStore();
const settingStore = useSettingStore();

const form = reactive<UpdateSettingRequest>({
  monitoring_interval_seconds: 10,
  alarm_check_interval_seconds: 60,
  alarm_repeat_interval_minutes: 15,
  system_cpu_critical_percent: 90.0,
  system_disk_critical_percent: 85.0,
  system_mem_critical_percent: 85.0,
  worker_heartbeat_timeout_seconds: 300,
});

const loadForm = () => {
  if (settingStore.settings) {
    form.monitoring_interval_seconds = settingStore.settings.monitoring_interval_seconds;
    form.alarm_check_interval_seconds = settingStore.settings.alarm_check_interval_seconds;
    form.alarm_repeat_interval_minutes = settingStore.settings.alarm_repeat_interval_minutes;
    form.system_cpu_critical_percent = settingStore.settings.system_cpu_critical_percent;
    form.system_disk_critical_percent = settingStore.settings.system_disk_critical_percent;
    form.system_mem_critical_percent = settingStore.settings.system_mem_critical_percent;
    form.worker_heartbeat_timeout_seconds = settingStore.settings.worker_heartbeat_timeout_seconds;
  }
};

onMounted(async () => {
  await settingStore.fetchSettings();
  loadForm();
});

watch(() => settingStore.settings, () => {
  loadForm();
}, { deep: true });

const saveSettings = async () => {
  await settingStore.updateSettings({ ...form });
};
</script>
