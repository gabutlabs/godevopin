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
