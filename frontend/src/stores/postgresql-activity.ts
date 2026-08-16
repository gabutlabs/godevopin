import type {
  PostgreSQLActivity,
  PostgreSQLLiveResponse,
  PostgreSQLTarget,
  PostgreSQLTargetRequest,
} from '@/types/postgresql-activity.type'
import { AxiosError } from 'axios'
import { defineStore } from 'pinia'
import axios from '@/plugins/axios'

export const usePostgreSQLActivityStore = defineStore('postgresqlActivity', {
  state: () => ({
    targets: [] as PostgreSQLTarget[],
    activities: [] as PostgreSQLActivity[],
    history: [] as PostgreSQLActivity[],
    loading: false,
    loadingTargets: false,
    loadingHistory: false,
    error: '',
    historyError: '',
  }),
  actions: {
    async fetchTargets () {
      this.loadingTargets = true
      try {
        const response = await axios.get('/postgresql-activity/targets')
        this.targets = response.data.data ?? []
        this.error = ''
      } catch (error) {
        this.error = this.errorMessage(error, 'Unable to load PostgreSQL targets.')
      } finally {
        this.loadingTargets = false
      }
      return this.targets
    },
    async fetchLive (targetId = 0, search = '') {
      this.loading = true
      try {
        const response = await axios.get('/postgresql-activity/live', {
          params: { target_id: targetId || undefined, search, limit: 500 },
        })
        const result = response.data.data as PostgreSQLLiveResponse
        this.activities = result.activities ?? []
        this.targets = result.targets ?? this.targets
        this.error = ''
      } catch (error) {
        this.activities = []
        this.error = this.errorMessage(error, 'Unable to load PostgreSQL activity.')
        console.error('Fetch PostgreSQL activity failed:', error)
      } finally {
        this.loading = false
      }
      return this.activities
    },
    async fetchHistory (activityKey: string, filter: string) {
      this.loadingHistory = true
      this.historyError = ''
      try {
        const response = await axios.get('/postgresql-activity/history', {
          params: { activity_key: activityKey, filter, limit: 500 },
        })
        this.history = response.data.data ?? []
      } catch (error) {
        this.history = []
        this.historyError = this.errorMessage(error, 'Unable to load PostgreSQL activity history.')
        console.error('Fetch PostgreSQL activity history failed:', error)
      } finally {
        this.loadingHistory = false
      }
      return this.history
    },
    async createTarget (payload: PostgreSQLTargetRequest) {
      try {
        await axios.post('/postgresql-activity/targets', payload)
        await this.fetchTargets()
        return { success: true, message: 'PostgreSQL target created successfully.' }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to create target.') }
      }
    },
    async updateTarget (id: number, payload: PostgreSQLTargetRequest) {
      try {
        await axios.put(`/postgresql-activity/targets/${id}`, payload)
        await this.fetchTargets()
        return { success: true, message: 'PostgreSQL target updated successfully.' }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to update target.') }
      }
    },
    async deleteTarget (id: number) {
      try {
        await axios.delete(`/postgresql-activity/targets/${id}`)
        this.targets = this.targets.filter(target => target.id !== id)
        this.activities = this.activities.filter(activity => activity.target_id !== id)
        return { success: true, message: 'PostgreSQL target deleted successfully.' }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'Unable to delete target.') }
      }
    },
    async testTarget (payload: PostgreSQLTargetRequest) {
      try {
        const response = await axios.post('/postgresql-activity/targets/test', payload)
        return {
          success: true,
          message: `Connection successful (PostgreSQL ${response.data.data.server_version}).`,
        }
      } catch (error) {
        return { success: false, message: this.errorMessage(error, 'PostgreSQL connection failed.') }
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
