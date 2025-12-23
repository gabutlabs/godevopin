import type { ContainerState } from "@/types/docker.type";

export function containerStateColor(state: ContainerState): string {
  let colorType = "";
  switch (true) {
    case state == "running":
      colorType = "green";
      break;
    case state == "exited":
      colorType = "red";
      break;
    case state == "dead":
      colorType = "red";
      break;
    case state == "paused":
      colorType = "yellow";
      break;
    default:
      colorType = "primary";
      break;
  }
  return colorType;
}
