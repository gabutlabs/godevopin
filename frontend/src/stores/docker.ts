import axios from "@/plugins/axios";
import type {
  ContainerDetailInfo,
  ContainerInfo,
  ImageInfo,
  NetworkInfo,
  VolumeInfo,
} from "@/types/docker.type";

export const useDockerStore = defineStore("docker", {
  state: () => ({
    containers: [] as Array<ContainerInfo>,
    container: {} as ContainerDetailInfo,
    images: [] as Array<ImageInfo>,
    networks: [] as Array<NetworkInfo>,
    volumes: [] as Array<VolumeInfo>,
    loading: false,
  }),
  actions: {
    async fetchAllContainers() {
      try {
        const response = await axios.get("/docker/containers");
        this.containers = response.data.data;
      } catch (error) {
        console.error("Fetch alarms failed:", error);
      }
    },
    async fetchContainerInfo(id: string) {
      this.loading = true;
      try {
        const response = await axios.get(`/docker/containers/${id}`);
        this.container = response.data.data;
        this.loading = false;
      } catch (error) {
        this.loading = false;
        console.error("Fetch alarms failed:", error);
      }
    },
    async fetchAllImages() {
      try {
        const response = await axios.get("/docker/images");
        this.images = response.data.data;
      } catch (error) {
        console.error("Fetch alarms failed:", error);
      }
    },
    async fetchAllNetworks() {
      try {
        const response = await axios.get("/docker/networks");
        this.networks = response.data.data;
      } catch (error) {
        console.error("Fetch alarms failed:", error);
      }
    },
    async fetchAllVolumes() {
      try {
        const response = await axios.get("/docker/volumes");
        this.volumes = response.data.data;
      } catch (error) {
        console.error("Fetch alarms failed:", error);
      }
    },
    async fetchAll() {
      try {
        await Promise.all([
          this.fetchAllImages(),
          this.fetchAllContainers(),
          this.fetchAllNetworks(),
          this.fetchAllVolumes(),
        ]);
      } catch (error) {
        console.error("error fetch all docker data: ", error);
      }
    },
  },
});
