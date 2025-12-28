<template>
  <div ref="terminalRef" class="terminal-container" />
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from "vue";
import { useRoute } from "vue-router";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import "@xterm/xterm/css/xterm.css";

interface ContainerExecProps {
  containerId?: string;
  command?: string;
}

const props = withDefaults(defineProps<ContainerExecProps>(), {
  containerId: undefined,
  command: "/bin/bash", // atau "/bin/sh"
});

const terminalRef = ref<HTMLElement | null>(null);
const term = ref<Terminal | null>(null);
const ws = ref<WebSocket | null>(null);
const fitAddon = ref<FitAddon | null>(null); // Ubah jadi ref
const isDisposed = ref(false); // Flag untuk cegah double dispose
const route = useRoute();

const routeContainerId = computed(
  () => (route.params as { id?: string }).id || ""
);
const containerId = computed(() => props.containerId || routeContainerId.value);

// WebSocket URL
const wsUrl = computed(() => {
  const id = containerId.value;
  const cmd = encodeURIComponent(props.command);
  return `ws://localhost:8080/ws/docker/container/exec/${id}?command=${cmd}`;
});

// Init WebSocket
const initWebSocket = () => {
  if (!term.value || isDisposed.value) return;

  try {
    ws.value = new WebSocket(wsUrl.value);

    ws.value.onopen = () => {
      console.log("WebSocket connected");
      if (term.value && !isDisposed.value) {
        term.value.clear();
        term.value.focus();

        // Kirim initial resize setelah connect
        setTimeout(() => {
          if (
            term.value &&
            ws.value?.readyState === WebSocket.OPEN &&
            !isDisposed.value
          ) {
            const payload = JSON.stringify({
              type: "resize",
              cols: term.value.cols,
              rows: term.value.rows,
            });
            ws.value.send(payload);

            // Trigger prompt
            setTimeout(() => {
              ws.value?.send("\r");
            }, 100);
          }
        }, 100);
      }
    };

    ws.value.onmessage = (event) => {
      if (!term.value || isDisposed.value) return;

      if (event.data) {
        if (event.data instanceof Blob) {
          const reader = new FileReader();
          reader.onload = () => {
            if (term.value && !isDisposed.value) {
              term.value.write(reader.result as string);
            }
          };
          reader.readAsText(event.data);
        } else if (typeof event.data === "string") {
          term.value.write(event.data);
        } else if (event.data instanceof ArrayBuffer) {
          const uint8Array = new Uint8Array(event.data);
          term.value.write(uint8Array);
        }
      }
    };

    ws.value.onerror = (error) => {
      console.error("WebSocket error:", error);
      if (term.value && !isDisposed.value) {
        term.value.writeln("\r\nConnection error\r\n");
      }
    };

    ws.value.onclose = (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      if (term.value && !isDisposed.value) {
        term.value.writeln("\r\nConnection closed\r\n");
      }
    };
  } catch (error) {
    console.error("Failed to create WebSocket:", error);
  }
};

// Handle terminal resize
const handleResize = () => {
  if (term.value && fitAddon.value && !isDisposed.value) {
    fitAddon.value.fit();

    // Kirim resize info ke backend
    if (ws.value?.readyState === WebSocket.OPEN) {
      const payload = JSON.stringify({
        type: "resize",
        cols: term.value.cols,
        rows: term.value.rows,
      });
      ws.value.send(payload);
    }
  }
};

// Init xterm
onMounted(() => {
  if (!terminalRef.value) return;

  try {
    term.value = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: "#1e1e1e",
        foreground: "#d4d4d4",
        cursor: "#ffffff",
        black: "#000000",
        red: "#cd3131",
        green: "#0dbc79",
        yellow: "#e5e510",
        blue: "#2472c8",
        magenta: "#bc3fbc",
        cyan: "#11a8cd",
        white: "#e5e5e5",
        brightBlack: "#666666",
        brightRed: "#f14c4c",
        brightGreen: "#23d18b",
        brightYellow: "#f5f543",
        brightBlue: "#3b8eea",
        brightMagenta: "#d670d6",
        brightCyan: "#29b8db",
        brightWhite: "#e5e5e5",
      },
      rows: 30,
      cols: 120,
    });

    fitAddon.value = new FitAddon();
    term.value.loadAddon(fitAddon.value);

    term.value.open(terminalRef.value);

    // Fit after a short delay to ensure proper sizing
    setTimeout(() => {
      if (fitAddon.value && !isDisposed.value) {
        fitAddon.value.fit();
      }
    }, 100);

    term.value.focus();

    // Handle user input
    term.value.onData((data) => {
      if (ws.value?.readyState === WebSocket.OPEN && !isDisposed.value) {
        ws.value.send(data);
      }
    });

    window.addEventListener("resize", handleResize);

    // Connect WebSocket
    initWebSocket();
  } catch (error) {
    console.error("Failed to initialize terminal:", error);
  }
});

// Cleanup
onBeforeUnmount(() => {
  if (isDisposed.value) return; //Cegah double cleanup

  isDisposed.value = true; //Set flag

  // Remove resize listener
  window.removeEventListener("resize", handleResize);

  // Close WebSocket
  if (ws.value) {
    try {
      ws.value.close();
    } catch (e) {
      console.error("Error closing WebSocket:", e);
    }
    ws.value = null;
  }

  // Dispose terminal (addon otomatis di-dispose)
  if (term.value) {
    try {
      term.value.dispose();
    } catch (e) {
      console.error("Error disposing terminal:", e);
    }
    term.value = null;
  }

  fitAddon.value = null;
});
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  min-height: 500px;
  background-color: #1e1e1e;
  padding: 8px;
  border-radius: 4px;
  overflow: hidden;
}

:deep(.xterm) {
  height: 100%;
  padding: 4px;
}

:deep(.xterm-viewport) {
  overflow-y: auto !important;
}

:deep(.xterm-screen) {
  height: 100% !important;
}
</style>
