import type {
  DatabaseActivity,
  DatabaseActivityLiveResponse,
  DatabaseEngine,
  DatabaseTarget,
  DatabaseTargetRequest,
} from '@/types/database-activity.type'
import { AxiosError } from 'axios'
import { defineStore } from 'pinia'
import axios from '@/plugins/axios'

const endpoints: Record<DatabaseEngine, string> = {
  postgresql: '/postgresql-activity',
  mysql: '/mysql-activity',
}

export const useDatabaseActivityStore = defineStore('databaseActivity', {
  state: () => ({
    postgresqlTargets: [] as DatabaseTarget[],
    mysqlTargets: [] as DatabaseTarget[],
    postgresqlActivities: [] as DatabaseActivity[],
    mysqlActivities: [] as DatabaseActivity[],
    history: [] as DatabaseActivity[],
    loading: false,
    loadingTargets: false,
    loadingHistory: false,
    error: '',
    historyError: '',
  }),
  getters: {
    targets: state => (engine: DatabaseEngine) => engine === 'postgresql' ? state.postgresqlTargets : state.mysqlTargets,
    activities: state => (engine: DatabaseEngine) => engine === 'postgresql' ? state.postgresqlActivities : state.mysqlActivities,
  },
  actions: {
    async fetchTargets () {
      this.loadingTargets = true
      try {
        const [postgresqlResponse, mysqlResponse] = await Promise.all([
          axios.get(`${endpoints.postgresql}/targets`),
          axios.get(`${endpoints.mysql}/targets`),
        ])
        this.postgresqlTargets = (postgresqlResponse.data.data ?? []).map((target: DatabaseTarget) => ({ ...target, engine: 'postgresql' as const }))
        this.mysqlTargets = (mysqlResponse.data.data ?? []).map((target: DatabaseTarget) => ({ ...target, engine: 'mysql' as const }))
        this.error = ''
      } catch (error) {
        this.error = this.errorMessage(error, 'Unable to load database targets.')
      } finally {
        this.loadingTargets = false
      }
    },
    async fetchLive (engine: DatabaseEngine, targetId = 0, search = '') {
      this.loading = true
      const normalizedTargetId = Number(targetId) || 0
      try {
        const response = await axios.get(`${endpoints[engine]}/live`, {
          params: { target_id: normalizedTargetId || undefined, search, limit: 500 },
        })
        const result = response.data.data as DatabaseActivityLiveResponse
        const activities = (result.activities ?? []).map(activity => ({ ...activity, engine }))
        const targets = (result.targets ?? []).map(target => ({ ...target, engine }))
        if (engine === 'postgresql') {
          this.postgresqlActivities = activities
          this.postgresqlTargets = targets
        } else {
          this.mysqlActivities = activities
          this.mysqlTargets = targets
        }
        this.error = ''
      } catch (error) {
        this.setActivities(engine, [])
        this.error = this.errorMessage(error, `Unable to load ${engine} activity.`)
        console.error(`Fetch ${engine} activity failed:`, error)
      } finally {
        this.loading = false
      }
      return this.activities(engine)
    },
    async fetchHistory (engine: DatabaseEngine, activityKey = '', filter = '1h', targetId = 0, search = '', blockedOnly = false) {
      this.loadingHistory = true
      this.historyError = ''
      const normalizedTargetId = Number(targetId) || 0
      try {
        const response = await axios.get(`${endpoints[engine]}/history`, {
          params: {
            activity_key: activityKey || undefined,
            target_id: normalizedTargetId || undefined,
            filter,
            search: search || undefined,
            blocked_only: blockedOnly || undefined,
            limit: 500,
          },
        })
        this.history = (response.data.data ?? [])
          .filter((activity: DatabaseActivity) => !normalizedTargetId || Number(activity.target_id) === normalizedTargetId)
          .map((activity: DatabaseActivity) => ({ ...activity, engine }))
      } catch (error) {
        this.history = []
        this.historyError = this.errorMessage(error, `Unable to load ${engine} activity history.`)
      } finally {
        this.loadingHistory = false
      }
      return this.history
    },
    async createTarget (engine: DatabaseEngine, payload: DatabaseTargetRequest) {
      try {
        await axios.post(`${endpoints[engine]}/targets`, payload)
        await this.fetchTargets()
        return { success: true, message: `${engineTitle(engine)} target created successfully.` }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to create target.') }
      }
    },
    async updateTarget (engine: DatabaseEngine, id: number, payload: DatabaseTargetRequest) {
      try {
        await axios.put(`${endpoints[engine]}/targets/${id}`, payload)
        await this.fetchTargets()
        return { success: true, message: `${engineTitle(engine)} target updated successfully.` }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to update target.') }
      }
    },
    async deleteTarget (engine: DatabaseEngine, id: number) {
      try {
        await axios.delete(`${endpoints[engine]}/targets/${id}`)
        this.setTargets(engine, this.targets(engine).filter(target => target.id !== id))
        this.setActivities(engine, this.activities(engine).filter(activity => activity.target_id !== id))
        return { success: true, message: `${engineTitle(engine)} target deleted successfully.` }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to delete target.') }
      }
    },
    async testTarget (engine: DatabaseEngine, payload: DatabaseTargetRequest) {
      try {
        const response = await axios.post(`${endpoints[engine]}/targets/test`, payload)
        return { success: true, message: `Connection successful (${engineTitle(engine)} ${response.data.data.server_version}).` }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, `${engineTitle(engine)} connection failed.`) }
      }
    },
    setTargets (engine: DatabaseEngine, targets: DatabaseTarget[]) {
      if (engine === 'postgresql') {
        this.postgresqlTargets = targets
      } else {
        this.mysqlTargets = targets
      }
    },
    setActivities (engine: DatabaseEngine, activities: DatabaseActivity[]) {
      if (engine === 'postgresql') {
        this.postgresqlActivities = activities
      } else {
        this.mysqlActivities = activities
      }
    },
    errorMessage (error: unknown, fallback: string) {
      if (error instanceof AxiosError) {
        return error.response?.data?.message || error.response?.data?.error || fallback
      }
      return fallback
    },
  },
})

function engineTitle (engine: DatabaseEngine) {
  return engine === 'postgresql' ? 'PostgreSQL' : 'MySQL'
}
