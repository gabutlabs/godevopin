<style scoped lang="scss"></style>
<template>
  <page-content title="Users">
    <template #append>
      <v-btn prepend-icon="mdi-plus" variant="outlined" @click="dialog = true"
        >Add User</v-btn
      >
    </template>
    <v-card>
      <v-card-text>
        <v-text-field
          v-model="name"
          class="ma-2"
          density="compact"
          placeholder="Search name..."
          width="30%"
        ></v-text-field>
        <v-data-table :items="state.users" :headers="headers">
          <template #item.actions="{ item }">
            <ActionTable :action-items="actionMenuItems" :item="item" />
          </template>
          <template #item.created_at="{ item }">
            {{ new Date(item.created_at).toLocaleString() }}
          </template>
          <template #item.updated_at="{ item }">
            {{ new Date(item.updated_at).toLocaleString() }}
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
    <v-dialog v-model="dialog" width="auto" max-width="700">
      <v-card
        max-width="700"
        width="700"
        prepend-icon="mdi-account"
        title="Form User"
      >
        <v-card-text>
          <v-form @submit.prevent="submit">
            <v-row>
              <v-col xl="6" md="6" sm="12">
                <v-text-field
                  v-model="payload.name"
                  label="Name"
                  :error-messages="v$.name.$errors.map((e: { $message: any; }) => e.$message)"
                  @blur="v$.name.$touch"
                  @input="v$.name.$touch"
                ></v-text-field>
              </v-col>
              <v-col xl="6" md="6" sm="12">
                <v-text-field
                  v-model="payload.email"
                  label="Email"
                  type="email"
                  :error-messages="v$.email.$errors.map((e: { $message: any; }) => e.$message)"
                  @blur="v$.email.$touch"
                  @input="v$.email.$touch"
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <p>Password default <b>Password123!</b></p>
              </v-col>
            </v-row>
            <v-btn class="mt-2" type="submit" block variant="elevated"
              >Submit</v-btn
            >
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </page-content>
</template>
<script setup lang="ts">
import { useNotify } from "@/composables/useNotify";
import { useUserStore, type User } from "@/stores/user";
import useVuelidate from "@vuelidate/core";
import { email, required } from "@vuelidate/validators";
import { ref, watch } from "vue";
const name = ref("");
const dialog = ref(false);
const state = useUserStore();
const payload = reactive({
  id: 0,
  email: "",
  name: "",
  is_updated: false,
});
const notify = useNotify();
const rules = {
  name: { required },
  email: { required, email },
};
const v$ = useVuelidate(rules, payload);
watch(name, async (newName) => {
  await state.setFilterUser(newName);
});
async function submit() {
  if (!(await v$.value.$validate())) {
    return;
  }
  if (payload.is_updated) {
    await state.updateUser(payload.id, payload);
  } else {
    await state.createUser(payload);
  }
  console.log(state.action_result);
  if (state.action_result.is_success) {
    // Login berhasil, arahkan ke halaman dashboard atau halaman yang diinginkan
    notify.success(state.action_result.message);
    dialog.value = false;
    payload.email = "";
    payload.name = "";
    await state.fetchAllUsers();
  } else {
    // Login gagal, tampilkan pesan kesalahan atau lakukan tindakan lain
    notify.error(
      `${state.action_result.message} : \n ${
        state.action_result.data != null
          ? Object.values(state.action_result.data).join("\n ")
          : ""
      }`
    );
  }
}
const actionMenuItems = [
  {
    title: "Edit",
    icon: "mdi-pencil",
    onClick: (item: User) => {
      payload.id = item.id;
      payload.name = item.name;
      payload.email = item.email;
      payload.is_updated = true;
      dialog.value = true;
    },
  },
  {
    title: "Delete",
    style: "color: red",
    icon: "mdi-delete",
    onClick: async (item: User) => {
      await state.deleteUser(item.id);
      if (state.action_result.is_success) {
        await state.fetchAllUsers();
        notify.success(state.action_result.message);
      } else {
        notify.error(state.action_result.message);
      }
    },
  },
];
const headers = [
  { title: "Name", value: "name" },
  { title: "Email", value: "email" },
  { title: "Created At", value: "created_at" },
  { title: "Updated At", value: "updated_at" },
  { title: "Actions", value: "actions", sortable: false },
];
onMounted(async () => {
  await state.fetchAllUsers();
});
</script>
