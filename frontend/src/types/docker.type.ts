export type ContainerState =
  | "created"
  | "running"
  | "paused"
  | "restarting"
  | "removing"
  | "exited"
  | "dead";
export interface ContainerInfo {
  id: string;
  name: string;
  image: string;
  status: string;
  state: ContainerState;
  labels: Record<string, string>;
}

export interface Port {
  IP?: string;
  PrivatePort: number; // uint16 → number
  PublicPort?: number;
  Type: string;
}

// ====== EndpointSettings ======
export interface EndpointSettings {
  IPAMConfig?: EndpointIPAMConfig | null;
  Links?: string[];
  Aliases?: string[];
  MacAddress?: string;
  DriverOpts?: Record<string, string>;
  GwPriority?: number;
  NetworkID?: string;
  EndpointID?: string;
  Gateway?: string;
  IPAddress?: string;
  IPPrefixLen?: number;
  IPv6Gateway?: string;
  GlobalIPv6Address?: string;
  GlobalIPv6PrefixLen?: number;
  DNSNames?: string[];
}

// Opsional: jika kamu pakai IPAMConfig
export interface EndpointIPAMConfig {
  // Sesuaikan jika kamu kirim field ini dari Go
  // Contoh umum:
  IPv4Address?: string;
  IPv6Address?: string;
  LinkLocalIPs?: string[];
}

// ====== NetworkSettingsSummary ======
export interface NetworkSettingsSummary {
  Networks?: Record<string, EndpointSettings | null>;
}

// ====== MountPoint ======
export interface MountPoint {
  Type?: string; // biasanya "bind", "volume", "tmpfs"
  Name?: string;
  Source: string;
  Destination: string;
  Driver?: string;
  Mode: string;
  RW: boolean;
  Propagation?: string; // biasanya "private", "rprivate", "slave", dll.
}

// ====== ContainerDetailInfo ======
export interface ContainerDetailInfo {
  container_info: ContainerInfo;
  ports: Port[];
  network_settings_summary?: NetworkSettingsSummary;
  mounts: MountPoint[];
}
export interface ImageInfo {
  id: string;
  tags: string[];
  size_mb: number; // atau gunakan sizeMb jika kamu ingin mengikuti gaya camelCase
  created: number; // Unix timestamp (dalam detik atau milidetik, sesuaikan)
}
export interface NetworkInfo {
  id: string;
  name: string;
  driver: string;
  scope: string;
}
export interface VolumeInfo {
  name: string;
  driver: string;
  mountpoint: string;
}
