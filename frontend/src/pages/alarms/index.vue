<template>
  <page-content title="Alarm">
    <template #append>
      <v-dialog max-width="500">
        <template v-slot:activator="{ props: activatorProps }">
          <v-btn
            prepend-icon="mdi-filter"
            variant="outlined"
            v-bind="activatorProps"
            >Filter</v-btn
          >
        </template>

        <template v-slot:default="{ isActive }">
          <v-card title="Dialog">
            <v-card-text>
              <v-combobox
                label="Select Status"
                :items="['FIRING', 'ACKNOWLEDGED']"
                variant="outlined"
                v-model="filterForm.status"
                clearable
              ></v-combobox>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>

              <v-btn text="Close" @click="isActive.value = false"></v-btn>
              <v-btn
                text="Apply"
                @click="applyFilter(() => (isActive.value = false))"
              ></v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </template>
    <v-tabs color="primary" hide-slider v-model="tab">
      <v-tab value="active">Active</v-tab>
      <v-tab value="history">History</v-tab>
    </v-tabs>

    <v-divider></v-divider>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="active">
        <ActiveAlarmComponent :filter-status="activeStatusFilter" />
      </v-tabs-window-item>
      <v-tabs-window-item value="history">
        <v-sheet class="pa-5">
          <HistoryAlarm :filter-status="historyStatusFilter" />
        </v-sheet>
      </v-tabs-window-item>
    </v-tabs-window>
  </page-content>
</template>
<script lang="ts" setup>
import PageContent from "@/components/PageContent.vue";
import ActiveAlarmComponent from "@/components/features/alarm/ActiveAlarm.vue";
const tab = ref("active");
const activeStatusFilter = ref("FIRING");
const historyStatusFilter = ref("");
const filterForm = reactive({
  status: "",
});
function applyFilter(cb: Function) {
  console.log("Applying filter:", filterForm.status);
  activeStatusFilter.value = filterForm.status || "FIRING";
  historyStatusFilter.value = filterForm.status || "";
  cb();
}
</script>
