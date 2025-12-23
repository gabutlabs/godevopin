<style scoped>
:deep(.v-overlay__content) {
  --v-overlay-anchor-origin: bottom right !important;
  transform-origin: right top !important;
  top: 60px !important;
  left: 1306px !important;
  min-width: 40px !important;
  max-width: 1446px !important;
  max-height: 403px !important;
}
</style>
<template>
  <v-layout>
    <v-app-bar elevation="1">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title text="Godevopin Admin" />
      <div class="d-flex gap-2 mx-4 align-center">
        <v-btn icon @click="() => router.push('/alarms')">
          <v-badge
            location="top right"
            color="warning"
            :content="alarmState.countActiveAlarms"
          >
            <v-icon icon="mdi-bell"></v-icon>
          </v-badge>
        </v-btn>
        <v-btn @click="state.changeTheme" icon>
          <v-icon>{{
            state.theme === "light" ? "mdi-weather-sunny" : "mdi-weather-night"
          }}</v-icon>
        </v-btn>
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-avatar color="red" v-bind="props">
              <span class="text-h5">CJ</span>
            </v-avatar>
          </template>
          <v-list>
            <template v-for="(item, index) in items">
              <v-list-item
                v-if="item.to"
                :key="index"
                :prepend-icon="item.icon"
                :to="item.to"
              >
                <v-list-item-title>{{ item.title }}</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-else
                :prepend-icon="item.icon"
                @click="item.onClick"
                :style="item.style"
              >
                <v-list-item-title>{{ item.title }}</v-list-item-title>
              </v-list-item>
            </template>
          </v-list>
        </v-menu>
      </div>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer">
      <v-list nav>
        <template v-for="v in sidebarMenu">
          <v-list-item
            :title="v.title"
            link
            :prepend-icon="v.icon"
            rounded
            slim
            :active="
              route.fullPath == v.link ||
              route.fullPath.startsWith(v.link + '/')
            "
            :to="v.link"
            color="#536DFE"
          >
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container class="overflow-y-auto" max-height="90vh">
        <router-view />
      </v-container>
    </v-main>
  </v-layout>
</template>

<script setup>
import sidebarMenu from "@/constants/sidebar-menu";
import { useAlarmStore } from "@/stores/alarm";
import { useAppStore } from "@/stores/app";
import { useAuthStore } from "@/stores/auth";
import { push } from "notivue";
import { ref } from "vue";
const drawer = ref(null);
const state = useAppStore();
const authState = useAuthStore();
const route = useRoute();
const router = useRouter();
const alarmState = useAlarmStore();
const items = ref([
  { title: "Profile", icon: "mdi-account", to: "/profile" },
  {
    title: "Sign-out",
    icon: "mdi-logout",
    onClick: () => {
      if (authState.logout()) {
        router.push("/login");
      } else {
        push.error("Failed logout, please try again");
      }
    },
    style: "color:#C62828",
  },
]);

onMounted(async () => {
  await alarmState.fetchCountActiveAlarm();
});
</script>
