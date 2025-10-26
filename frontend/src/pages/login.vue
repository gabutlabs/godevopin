<route lang="yaml">
meta:
  layout: auth-layout
  guestOnly: true
</route>
<style scoped>
a {
  text-decoration: none;
}
</style>
<template>
  <v-card class="pa-4">
    <v-card-text class="pt-2">
      <h4 class="text-h4 mb-1">Welcome to Vtify! 👋🏻</h4>
      <p class="mb-0">Please sign-in to your account and start the adventure</p>
    </v-card-text>

    <v-card-text>
      <v-form @submit.prevent="submit">
        <v-row>
          <v-col cols="12">
            <v-text-field
              label="Email"
              class="mb-2"
              v-model="formState.email"
              :error-messages="v$.email.$errors.map((e: { $message: any; }) => e.$message)"
              @blur="v$.email.$touch"
              @input="v$.email.$touch"
            ></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-text-field
              label="Password"
              :type="isPasswordVisible ? 'text' : 'password'"
              autocomplete="password"
              :append-inner-icon="isPasswordVisible ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="isPasswordVisible = !isPasswordVisible"
              :error-messages="v$.password.$errors.map((e: { $message: any; }) => e.$message)"
              v-model="formState.password"
              @blur="v$.password.$touch"
              @input="v$.password.$touch"
            ></v-text-field>
            <div class="d-flex justify-end my-6">
              <router-link class="text-primary" to="/forgot-password">
                Forgot Password?
              </router-link>
            </div>
            <v-btn
              block
              color="primary"
              class="mt-4"
              rounded="lg"
              variant="elevated"
              density="default"
              type="submit"
              :loading="loading"
              >Masuk</v-btn
            >
          </v-col>
          <v-col cols="12" class="text-center text-base">
            <span>New on our platform?</span>
            <router-link class="text-primary ms-2" to="/register">
              Create an account
            </router-link>
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { useNotify } from "@/composables/useNotify";
import { useAuthStore } from "@/stores/auth";
import useVuelidate from "@vuelidate/core";
import { email, required } from "@vuelidate/validators";
const isPasswordVisible = ref(false);
// Script untuk halaman login Anda
const notify = useNotify();
const formState = reactive({
  email: "",
  password: "",
});
const rules = {
  password: { required },
  email: { required, email },
};
const v$ = useVuelidate(rules, formState);
const state = useAuthStore();
const router = useRouter();
const loading = ref(false);
async function submit() {
  if (!(await v$.value.$validate())) {
    return;
  }
  loading.value = true;
  const result = await state.login(formState);
  loading.value = false;
  if (result) {
    // Login berhasil, arahkan ke halaman dashboard atau halaman yang diinginkan
    router.push("/"); // Ganti dengan rute yang sesuai
  } else {
    console.error("hai rerp");
    // Login gagal, tampilkan pesan kesalahan atau lakukan tindakan lain
    notify.error("Login failed. Please check your credentials.");
  }
}
</script>
