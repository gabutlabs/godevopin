<template>
  <v-menu>
    <template v-slot:activator="{ props }">
      <v-btn
        color="primary"
        icon="mdi-dots-horizontal"
        v-bind="props"
        density="compact"
        variant="text"
      />
    </template>
    <v-list>
      <v-list-item
        v-for="(action, index) in actionItems"
        :key="index"
        :value="index"
        :disabled="
          typeof action.disabled === 'function'
            ? action.disabled(item)
            : action.disabled
        "
        @click="action.onClick(item)"
      >
        <v-list-item-title
          ><v-icon :icon="action.icon" :style="action.style"></v-icon>
          {{ action.title }}</v-list-item-title
        >
      </v-list-item>
    </v-list>
  </v-menu>
</template>
<script lang="ts" setup>
type ActionItem<T = any> = {
  title: string;
  disabled?: boolean | ((item: T) => boolean);
  onClick: (item: T) => void;
  icon?: string;
  style?: string;
};

defineProps<{
  actionItems: ActionItem[];
  item?: any;
}>();
</script>
