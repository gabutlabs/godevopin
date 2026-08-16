export interface PostgreSQLTarget {
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
}

export interface PostgreSQLTargetRequest {
  name: string
  host: string
  port: number
  database: string
  username: string
  password?: string
  ssl_mode: string
  enabled: boolean
}

export interface PostgreSQLActivity {
  activity_key: string
  target_id: number
  target_name: string
  server_version: string
  pid: number
  database_name: string
  username: string
  application_name?: string
  client_address?: string
  client_port?: number
  backend_type?: string
  backend_start?: string
  transaction_start?: string
  query_start?: string
  state_change?: string
  wait_event_type?: string
  wait_event?: string
  blocking_pid?: number
  state: string
  query?: string
  query_fingerprint?: string
  query_duration_ms: number
  observed_at: string
}

export interface PostgreSQLLiveResponse {
  activities: PostgreSQLActivity[]
  targets: PostgreSQLTarget[]
}
