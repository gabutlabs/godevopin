export interface ProcessMetric {
  process_key: string;
  pid: number;
  process_start_time: number;
  parent_pid: number;
  name: string;
  executable?: string;
  username?: string;
  status: string;
  cpu_percent: number;
  memory_percent: number;
  memory_bytes: number;
  virtual_memory_bytes: number;
  thread_count: number;
  observed_at: string;
}

export interface ProcessHistoryPoint {
  time_interval: string;
  avg_cpu_percent: number;
  avg_memory_percent: number;
  avg_memory_bytes: number;
  max_cpu_percent: number;
}
