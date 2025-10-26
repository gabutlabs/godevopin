import { ref } from "vue";

interface NotifyOptions {
  color?: string;
  timeout?: number;
  icon?: string;
}

const snackbar = ref(false);
const message = ref("");
const color = ref("info");
const timeout = ref(3000);
const icon = ref("");

export function useNotify() {
  const notify = (msg: string, options: NotifyOptions = {}) => {
    message.value = msg;
    color.value = options.color || "info";
    timeout.value = options.timeout || 3000;
    icon.value = options.icon || "";
    snackbar.value = true;
  };

  const success = (msg: string, options: NotifyOptions = {}) => {
    notify(msg, { ...options, color: "success", icon: "mdi-check-circle" });
  };

  const error = (msg: string, options: NotifyOptions = {}) => {
    notify(msg, { ...options, color: "error", icon: "mdi-alert-circle" });
  };

  const warning = (msg: string, options: NotifyOptions = {}) => {
    notify(msg, { ...options, color: "warning", icon: "mdi-alert" });
  };

  const info = (msg: string, options: NotifyOptions = {}) => {
    notify(msg, { ...options, color: "info", icon: "mdi-information" });
  };

  return {
    snackbar,
    message,
    color,
    timeout,
    icon,
    notify,
    success,
    error,
    warning,
    info,
  };
}
