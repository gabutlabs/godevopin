<template>
  <div class="w-full d-flex flex-column">
    <div class="d-flex justify-end mb-4">
      <v-btn
        :color="status === 'OPEN' ? 'success' : 'error'"
        :disabled="!containerId"
        @click="toggleConnection"
        :prepend-icon="status === 'OPEN' ? 'mdi-stop' : 'mdi-play'"
      >
        {{ status === "OPEN" ? "Connected" : "Connect" }}
        <v-tooltip activator="parent" location="top">
          {{
            status === "OPEN"
              ? "Click to disconnect from logs"
              : "Click to connect to real-time logs"
          }}
        </v-tooltip>
      </v-btn>
    </div>

    <div class="w-full border rounded-lg d-flex flex-column pa-3">
      <div class="d-flex justify-space-between mb-2">
        <h4>Container Logs</h4>
        <v-btn
          size="small"
          variant="text"
          icon="mdi-content-copy"
          @click="copyLogs"
          :disabled="logs.length === 0"
        >
          <v-tooltip activator="parent" location="top"
            >Copy logs to clipboard</v-tooltip
          >
        </v-btn>
      </div>

      <div
        ref="logsContainerRef"
        class="logs-container pa-3 font-mono text-sm"
        style="
          background-color: #1e1e1e;
          color: #d4d4d4;
          max-height: 500px;
          overflow-y: auto;
          font-family: 'Courier New', monospace;
        "
      >
        <div
          v-for="(log, index) in logs"
          :key="index"
          class="log-line py-1"
          :class="{ 'error-line': log.includes('[STDERR]') }"
        >
          {{ log }}
        </div>
        <div v-if="logs.length === 0" class="text-center py-10 text-gray">
          No logs available. Connect to start receiving real-time logs.
        </div>
      </div>

      <div class="d-flex justify-space-between align-center mt-3">
        <div class="text-caption">
          {{ logs.length }} log{{ logs.length !== 1 ? "s" : "" }} shown
        </div>
        <v-btn
          size="small"
          variant="text"
          icon="mdi-trash-can"
          @click="clearLogs"
          :disabled="logs.length === 0"
        >
          <v-tooltip activator="parent" location="top">Clear logs</v-tooltip>
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import { useWebSocket } from "@vueuse/core";
import { useNotify } from "@/composables/useNotify";

interface ContainerLogsProps {
  containerId?: string;
}

const props = withDefaults(defineProps<ContainerLogsProps>(), {
  containerId: undefined,
});

// Get container ID from route if not provided as prop
const route = useRoute();
const routeContainerId = computed(
  () => (route.params as { id?: string }).id || ""
);
const containerId = computed(() => props.containerId || routeContainerId.value);

// State variables
const logs = ref<string[]>([]);
const logsContainerRef = ref<HTMLElement | null>(null);
const notify = useNotify();

// WebSocket URL computed property
const wsUrl = computed(() => {
  if (!containerId.value) return undefined;
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//localhost:8080/ws/docker/container/logs/${containerId.value}`;
});

// UseWebSocket composable
const { status, data, close, open } = useWebSocket(wsUrl, {
  autoReconnect: false,
  immediate: false, // Don't connect automatically
  onConnected: () => {
    console.log("Connected to container logs WebSocket");
    logs.value = [];
    notify.success(
      `Connected to logs for container ${containerId.value.substring(0, 12)}`
    );
  },
  onDisconnected: () => {
    console.log("Disconnected from container logs WebSocket");
    notify.warning(
      `Disconnected from logs for container ${containerId.value.substring(
        0,
        12
      )}`
    );
  },
  onError: (event) => {
    console.error("WebSocket error:", event);
    notify.error("WebSocket connection error");
  },
});

// Watch for incoming WebSocket data
watch(data, (newData) => {
  if (newData) {
    addLogLine(newData);
  }
});

// Function to toggle connection
const toggleConnection = () => {
  if (status.value === "OPEN") {
    close();
    status.value = "CLOSED";
  } else {
    if (wsUrl.value) {
      open();
    }
  }
};

// Function to add a log line to the display
const addLogLine = (logLine: string) => {
  // Only add non-empty lines
  if (logLine.trim()) {
    logs.value.push(logLine);

    // Limit the number of logs to prevent memory issues
    if (logs.value.length > 1000) {
      logs.value.shift(); // Remove the oldest log
    }

    // Scroll to bottom after next render
    nextTick(() => {
      scrollToBottom();
    });
  }
};

// Function to scroll to the bottom of the logs container
const scrollToBottom = () => {
  if (logsContainerRef.value) {
    logsContainerRef.value.scrollTop = logsContainerRef.value.scrollHeight;
  }
};

// Function to clear all logs
const clearLogs = () => {
  logs.value = [];
  notify.info("Logs cleared");
};

// Function to copy logs to clipboard
const copyLogs = async () => {
  try {
    await navigator.clipboard.writeText(logs.value.join("\n"));
    notify.success("Logs copied to clipboard");
  } catch (err) {
    console.error("Failed to copy logs:", err);
    notify.error("Failed to copy logs to clipboard");
  }
};

// Watch for changes in container ID and reconnect if needed
watch(containerId, (newId, oldId) => {
  if (newId && newId !== oldId) {
    if (status.value === "OPEN") {
      // Close current connection
      close();
      // Clear logs for new container
      logs.value = [];
      // Open new connection after a short delay
      nextTick(() => {
        setTimeout(() => {
          open();
        }, 100);
      });
    }
  }
});

// Cleanup on component unmount
onUnmounted(() => {
  if (status.value === "OPEN") {
    close();
  }
});
</script>

<style scoped>
.logs-container {
  background-color: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  font-family: "Courier New", monospace;
  font-size: 12px;
  line-height: 1.4;
}

.error-line {
  color: #f44336; /* Red color for error lines */
}

.log-line {
  white-space: pre-wrap; /* Preserve whitespace and wrap lines */
  word-break: break-all; /* Break long lines */
}
</style>
