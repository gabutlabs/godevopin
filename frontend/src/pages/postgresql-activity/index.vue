<route lang="yaml">
meta:
  requiresAuth: true
</route>

<template>
  <PageContent title="PostgreSQL Activity">
    <template #append>
      <div class="d-flex flex-wrap ga-3 mx-4 align-center">
        <v-select
          v-model="selectedTargetId"
          density="compact"
          hide-details
          :items="targetItems"
          label="Target"
          style="min-width: 190px"
          variant="outlined"
        />
        <v-text-field
          v-model="search"
          clearable
          density="compact"
          hide-details
          label="Search activity"
          prepend-inner-icon="mdi-magnify"
          style="min-width: 220px"
          variant="outlined"
        />
        <v-select
          v-model="stateFilter"
          density="compact"
          hide-details
          :items="stateItems"
          label="State"
          style="width: 160px"
          variant="outlined"
        />
        <v-select
          v-model="refreshSeconds"
          density="compact"
          hide-details
          :items="refreshItems"
          label="Refresh"
          style="width: 130px"
          variant="outlined"
        />
        <v-btn
          :prepend-icon="isRefreshing ? 'mdi-pause' : 'mdi-play'"
          variant="outlined"
          @click="toggleRefresh"
        >
          {{ isRefreshing ? "Pause" : "Resume" }}
        </v-btn>
        <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreateTarget">
          Add target
        </v-btn>
      </div>
    </template>

    <v-alert v-if="activityStore.error" class="mb-4" closable type="error">
      {{ activityStore.error }}
    </v-alert>

    <v-alert v-if="activityStore.targets.length === 0" class="mb-4" type="info">
      Add a PostgreSQL target to start monitoring connections and queries.
    </v-alert>

    <v-row>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">ACTIVE SESSIONS</div>
            <div class="text-h4">{{ visibleActivities.length }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">WAITING OR BLOCKED</div>
            <div class="text-h4">{{ waitingActivities }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline">LONGEST QUERY</div>
            <div class="text-h4">{{ formatDuration(longestActivity) }}</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="d-flex align-center">
            <span>Live connections and queries</span>
            <v-spacer />
            <v-progress-circular v-if="activityStore.loading" indeterminate size="20" width="2" />
          </v-card-title>
          <v-data-table
            :headers="activityHeaders"
            item-value="activity_key"
            :items="visibleActivities"
            :items-per-page="25"
            :loading="activityStore.loading"
          >
            <template #item.target_name="{ item }">
              <div class="font-weight-medium">{{ item.target_name }}</div>
              <div class="text-caption text-medium-emphasis">PID {{ item.pid }}</div>
            </template>
            <template #item.database_name="{ item }">
              <div>{{ item.database_name || "-" }}</div>
              <div class="text-caption text-medium-emphasis">{{ item.username || "unknown user" }}</div>
            </template>
            <template #item.state="{ item }">
              <v-chip :color="stateColor(item)" size="small">
                {{ item.state || "unknown" }}
              </v-chip>
              <div v-if="item.wait_event" class="text-caption text-medium-emphasis">
                {{ item.wait_event_type ? `${item.wait_event_type}: ` : "" }}{{ item.wait_event }}
              </div>
              <div v-if="item.blocking_pid" class="text-caption text-error">
                Blocking PID {{ item.blocking_pid }}
              </div>
            </template>
            <template #item.query_duration_ms="{ item }">
              <span :class="item.query_duration_ms >= 5000 ? 'text-warning font-weight-medium' : ''">
                {{ formatDuration(item.query_duration_ms) }}
              </span>
            </template>
            <template #item.query="{ item }">
              <div class="text-truncate" style="max-width: 520px" :title="item.query">
                {{ item.query || "-" }}
              </div>
            </template>
            <template #item.actions="{ item }">
              <v-tooltip location="top" text="View activity history">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    aria-label="View activity history"
                    icon="mdi-chart-line"
                    size="small"
                    variant="text"
                    @click.stop="openHistory(item)"
                  />
                </template>
              </v-tooltip>
            </template>
          </v-data-table>
        </v-card>
      </v-col>

      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title>Configured targets</v-card-title>
          <v-data-table
            :headers="targetHeaders"
            item-value="id"
            :items="activityStore.targets"
            :items-per-page="10"
            :loading="activityStore.loadingTargets"
          >
            <template #item.enabled="{ item }">
              <v-chip :color="item.enabled ? 'success' : 'secondary'" size="small">
                {{ item.enabled ? "enabled" : "disabled" }}
              </v-chip>
            </template>
            <template #item.last_checked_at="{ item }">
              {{ item.last_checked_at ? new Date(item.last_checked_at).toLocaleString() : "Never" }}
            </template>
            <template #item.last_error="{ item }">
              <span v-if="item.last_error" class="text-error text-caption">{{ item.last_error }}</span>
              <span v-else class="text-success text-caption">Healthy</span>
            </template>
            <template #item.actions="{ item }">
              <v-btn
                aria-label="Edit target"
                icon="mdi-pencil"
                size="small"
                variant="text"
                @click="openEditTarget(item)"
              />
              <v-btn
                aria-label="Delete target"
                color="error"
                icon="mdi-delete"
                size="small"
                variant="text"
                @click="removeTarget(item)"
              />
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="targetDialog" max-width="720">
      <v-card>
        <v-card-title>{{ editingTargetId ? "Edit PostgreSQL target" : "Add PostgreSQL target" }}</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="saveTarget">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field v-model="targetForm.name" label="Name" required variant="outlined" />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="targetForm.host" label="Host" required variant="outlined" />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="targetForm.port"
                  label="Port"
                  required
                  type="number"
                  variant="outlined"
                />
              </v-col>
              <v-col cols="12" sm="8">
                <v-text-field v-model="targetForm.database" label="Database" required variant="outlined" />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="targetForm.username" label="Monitoring username" required variant="outlined" />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="targetForm.password"
                  :hint="editingTargetId ? 'Leave blank to keep the current password.' : ''"
                  label="Password"
                  persistent-hint
                  :required="!editingTargetId"
                  type="password"
                  variant="outlined"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-select v-model="targetForm.ssl_mode" :items="sslModes" label="SSL mode" variant="outlined" />
              </v-col>
              <v-col class="d-flex align-center" cols="12" sm="6">
                <v-switch v-model="targetForm.enabled" color="primary" hide-details label="Collect activity" />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-btn variant="text" @click="targetDialog = false">Cancel</v-btn>
          <v-spacer />
          <v-btn :disabled="!targetFormReady" :loading="testingTarget" variant="outlined" @click="testTarget">
            Test connection
          </v-btn>
          <v-btn color="primary" :disabled="!targetFormReady" :loading="savingTarget" @click="saveTarget">
            Save target
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="historyDialog" max-width="1200">
      <v-card>
        <v-card-title class="d-flex align-center">
          <span>Activity history</span>
          <v-spacer />
          <v-select
            v-model="historyFilter"
            density="compact"
            hide-details
            :items="historyItems"
            style="width: 150px"
            variant="outlined"
          />
          <v-btn icon="mdi-close" variant="text" @click="historyDialog = false" />
        </v-card-title>
        <v-card-subtitle v-if="selectedActivity" class="pb-2">
          {{ selectedActivity.target_name }} / PID {{ selectedActivity.pid }} / {{ selectedActivity.database_name }}
        </v-card-subtitle>
        <v-card-text>
          <div v-if="activityStore.loadingHistory" class="text-center py-8">
            <v-progress-circular color="primary" indeterminate />
            <div class="mt-3">Loading activity history...</div>
          </div>
          <div v-else-if="activityStore.historyError" class="text-center text-error py-8">
            {{ activityStore.historyError }}
          </div>
          <v-data-table v-else :headers="historyHeaders" :items="activityStore.history" :items-per-page="15">
            <template #item.observed_at="{ item }">{{ new Date(item.observed_at).toLocaleString() }}</template>
            <template #item.query_duration_ms="{ item }">{{ formatDuration(item.query_duration_ms) }}</template>
            <template #item.state="{ item }"><v-chip :color="stateColor(item)" size="small">{{ item.state }}</v-chip></template>
            <template #item.query="{ item }"><div class="text-wrap" style="max-width: 650px">{{ item.query || "-" }}</div></template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-dialog>
  </PageContent>
</template>

<script lang="ts" setup>
  import type {
    PostgreSQLActivity,
    PostgreSQLTarget,
    PostgreSQLTargetRequest,
  } from '@/types/postgresql-activity.type'
  import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
  import PageContent from '@/components/PageContent.vue'
  import { useNotify } from '@/composables/useNotify'
  import { usePostgreSQLActivityStore } from '@/stores/postgresql-activity'

  const notify = useNotify()
  const activityStore = usePostgreSQLActivityStore()
  const selectedTargetId = ref(0)
  const search = ref('')
  const stateFilter = ref('all')
  const refreshSeconds = ref(10)
  const isRefreshing = ref(true)
  const targetDialog = ref(false)
  const editingTargetId = ref(0)
  const savingTarget = ref(false)
  const testingTarget = ref(false)
  const historyDialog = ref(false)
  const historyFilter = ref('1h')
  const selectedActivity = ref<PostgreSQLActivity | null>(null)
  let refreshTimer: number | null = null

  const targetForm = reactive<PostgreSQLTargetRequest>({
    name: '',
    host: '',
    port: 5432,
    database: 'postgres',
    username: '',
    password: '',
    ssl_mode: 'prefer',
    enabled: true,
  })

  const refreshItems = [
    { title: '5 seconds', value: 5 },
    { title: '10 seconds', value: 10 },
    { title: '30 seconds', value: 30 },
  ]
  const stateItems = [
    { title: 'All states', value: 'all' },
    { title: 'Active', value: 'active' },
    { title: 'Idle in transaction', value: 'idle in transaction' },
  ]
  const historyItems = [
    { title: 'Last 1 hour', value: '1h' },
    { title: 'Last 6 hours', value: '6h' },
    { title: 'Last 1 day', value: '1d' },
    { title: 'Last 7 days', value: '7d' },
    { title: 'Last 30 days', value: '30d' },
  ]
  const sslModes = ['prefer', 'require', 'verify-ca', 'verify-full', 'disable']
  const activityHeaders = [
    { title: 'Target / PID', key: 'target_name' },
    { title: 'Database / User', key: 'database_name' },
    { title: 'State / Wait', key: 'state' },
    { title: 'Duration', key: 'query_duration_ms', align: 'end' as const },
    { title: 'Query', key: 'query' },
    { title: '', key: 'actions', sortable: false, align: 'end' as const },
  ]
  const targetHeaders = [
    { title: 'Name', key: 'name' },
    { title: 'Host', key: 'host' },
    { title: 'Database', key: 'database' },
    { title: 'SSL', key: 'ssl_mode' },
    { title: 'Status', key: 'enabled' },
    { title: 'Last check', key: 'last_checked_at' },
    { title: 'Health', key: 'last_error' },
    { title: '', key: 'actions', sortable: false, align: 'end' as const },
  ]
  const historyHeaders = [
    { title: 'Observed', key: 'observed_at' },
    { title: 'State', key: 'state' },
    { title: 'Duration', key: 'query_duration_ms' },
    { title: 'Wait', key: 'wait_event' },
    { title: 'Query', key: 'query' },
  ]

  const targetItems = computed(() => [
    { title: 'All enabled targets', value: 0 },
    ...activityStore.targets.map(target => ({ title: target.name, value: target.id })),
  ])
  const visibleActivities = computed(() => {
    if (stateFilter.value === 'all') return activityStore.activities
    return activityStore.activities.filter(activity => activity.state === stateFilter.value)
  })
  const waitingActivities = computed(() =>
    visibleActivities.value.filter(activity => Boolean(activity.wait_event)).length,
  )
  const longestActivity = computed(() =>
    visibleActivities.value.reduce((max, activity) => Math.max(max, activity.query_duration_ms), 0),
  )
  const targetFormReady = computed(() =>
    Boolean(targetForm.name.trim() && targetForm.host.trim() && targetForm.username.trim() && targetForm.database.trim() && (editingTargetId.value || targetForm.password)),
  )

  function formatDuration (milliseconds: number) {
    const value = Number(milliseconds || 0)
    if (value < 1000) return `${value} ms`
    if (value < 60_000) return `${(value / 1000).toFixed(1)} s`
    return `${Math.floor(value / 60_000)}m ${Math.floor((value % 60_000) / 1000)}s`
  }

  function stateColor (activity: PostgreSQLActivity) {
    if (activity.wait_event) return 'warning'
    if (activity.state === 'active') return 'success'
    if (activity.state.includes('aborted')) return 'error'
    return 'secondary'
  }

  async function fetchLive () {
    await activityStore.fetchLive(selectedTargetId.value, search.value)
  }

  function stopRefreshTimer () {
    if (refreshTimer !== null) {
      window.clearInterval(refreshTimer)
      refreshTimer = null
    }
  }

  function startRefreshTimer () {
    stopRefreshTimer()
    if (isRefreshing.value) refreshTimer = window.setInterval(fetchLive, refreshSeconds.value * 1000)
  }

  function toggleRefresh () {
    isRefreshing.value = !isRefreshing.value
    startRefreshTimer()
  }

  function resetTargetForm () {
    Object.assign(targetForm, {
      name: '',
      host: '',
      port: 5432,
      database: 'postgres',
      username: '',
      password: '',
      ssl_mode: 'prefer',
      enabled: true,
    })
  }

  function openCreateTarget () {
    editingTargetId.value = 0
    resetTargetForm()
    targetDialog.value = true
  }

  function openEditTarget (target: PostgreSQLTarget) {
    editingTargetId.value = target.id
    Object.assign(targetForm, {
      name: target.name,
      host: target.host,
      port: target.port,
      database: target.database,
      username: target.username,
      password: '',
      ssl_mode: target.ssl_mode,
      enabled: target.enabled,
    })
    targetDialog.value = true
  }

  async function testTarget () {
    testingTarget.value = true
    const result = await activityStore.testTarget({ ...targetForm })
    testingTarget.value = false
    result.success ? notify.success(result.message) : notify.error(result.message)
  }

  async function saveTarget () {
    savingTarget.value = true
    const result = editingTargetId.value
      ? await activityStore.updateTarget(editingTargetId.value, { ...targetForm })
      : await activityStore.createTarget({ ...targetForm })
    savingTarget.value = false
    if (!result.success) {
      notify.error(result.message)
      return
    }
    targetDialog.value = false
    notify.success(result.message)
    await fetchLive()
  }

  async function removeTarget (target: PostgreSQLTarget) {
    if (!window.confirm(`Delete PostgreSQL target "${target.name}"?`)) return
    const result = await activityStore.deleteTarget(target.id)
    result.success ? notify.success(result.message) : notify.error(result.message)
    if (selectedTargetId.value === target.id) selectedTargetId.value = 0
  }

  async function openHistory (activity: PostgreSQLActivity) {
    selectedActivity.value = activity
    historyDialog.value = true
    await activityStore.fetchHistory(activity.activity_key, historyFilter.value)
  }

  watch([selectedTargetId, search], fetchLive)
  watch([refreshSeconds, isRefreshing], startRefreshTimer)
  watch(historyFilter, async filter => {
    if (selectedActivity.value) await activityStore.fetchHistory(selectedActivity.value.activity_key, filter)
  })

  onMounted(async () => {
    await activityStore.fetchTargets()
    await fetchLive()
    startRefreshTimer()
  })

  onUnmounted(stopRefreshTimer)
</script>
