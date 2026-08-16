export type DatabaseEngine = 'postgresql' | 'mysql'

export interface DatabaseTarget {
  id: number
  name: string
  host: string
  port: number
  database: string
  username: string
  ssl_mode: string
  enabled: boolean
  password_configured: boolean
  last_checked_at?: string
  last_error?: string
  created_at: string
  updated_at: string
  engine: DatabaseEngine
}

export interface DatabaseTargetRequest {
  name: string
  host: string
  port: number
  database: string
  username: string
  password?: string
  ssl_mode: string
  enabled: boolean
}

export interface DatabaseActivity {
  activity_key: string
  target_id: number
  target_name: string
  server_version: string
  pid: number
  database_name: string
  username: string
  application_name?: string
  client_address?: string
  backend_type?: string
  command?: string
  wait_event_type?: string
  wait_event?: string
  blocking_pid?: number
  blocking_database?: string
  blocking_username?: string
  blocking_state?: string
  blocking_query?: string
  state: string
  query?: string
  query_fingerprint?: string
  query_duration_ms: number
  observed_at: string
  engine: DatabaseEngine
}

export interface DatabaseActivityLiveResponse {
  activities: DatabaseActivity[]
  targets: DatabaseTarget[]
}
