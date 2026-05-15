/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";

// Composables
import { createVuetify } from "vuetify";
import type { ThemeDefinition } from "vuetify";
import { el } from "vuetify/locale";
/**
 * TEMA LIGHT
 * Background abu-abu sangat muda dengan card/permukaan berwarna putih bersih.
 * Memberikan kesan lapang dan bersih.
 */
const lightTheme: ThemeDefinition = {
  dark: false,
  colors: {
    background: "#F9FAFB", // Warna latar belakang utama yang soft
    surface: "#FFFFFF", // Warna untuk card, sheet, menu, dll.
    primary: "#536DFE", // Warna utama sesuai permintaan Anda
    "primary-darken-1": "#3D5AFE", // Versi lebih gelap untuk hover/active
    secondary: "#00BFA5", // Warna sekunder (teal) yang cocok dengan primary
    "secondary-darken-1": "#00897B",
    error: "#E53935", // Merah untuk error
    info: "#2196F3", // Biru untuk info
    success: "#43A047", // Hijau untuk sukses
    warning: "#FB8C00", // Oranye untuk peringatan

    // Warna teks & ikon
    "on-surface": "#424242",
    "on-background": "#424242",
    "on-primary": "#FFFFFF",
    "on-secondary": "#FFFFFF",
  },
  variables: {
    "border-color": "#000000",
    "border-opacity": 0.12,
    "high-emphasis-opacity": 0.87,
    "medium-emphasis-opacity": 0.6,
    "disabled-opacity": 0.38,
    "idle-opacity": 0.04,
    "hover-opacity": 0.04,
    "focus-opacity": 0.12,
    "selected-opacity": 0.08,
    "activated-opacity": 0.12,
    "pressed-opacity": 0.12,
    "dragged-opacity": 0.08,
    "theme-kbd": "#212529",
    "theme-on-kbd": "#FFFFFF",
    "theme-code": "#F5F5F5",
    "theme-on-code": "#000000",
  },
};

/**
 * TEMA DARK
 * Menggunakan prinsip Material Design, di mana background tidak hitam pekat
 * dan surface sedikit lebih terang untuk menciptakan efek kedalaman (depth).
 */
const darkTheme: ThemeDefinition = {
  dark: true,
  colors: {
    background: "#0b1326",
    surface: "#0b1326",
    "surface-bright": "#31394d",
    "surface-variant": "#2d3449",
    primary: "#adc6ff",
    "primary-darken-1": "#4d8eff", // Using primary-container
    secondary: "#b7c8e1",
    "secondary-darken-1": "#3a4a5f", // Using secondary-container
    error: "#ffb4ab",
    info: "#adc6ff",
    success: "#b7c8e1",
    warning: "#ffb786",

    // Text & Icon Colors
    "on-surface": "#dae2fd",
    "on-surface-variant": "#c2c6d6",
    "on-background": "#dae2fd",
    "on-primary": "#002e6a",
    "on-secondary": "#213145",
    "on-error": "#690005",
    "on-warning": "#502400",
  },
  variables: {
    "border-color": "#FFFFFF",
    "border-opacity": 0.12,
    "high-emphasis-opacity": 0.87,
    "medium-emphasis-opacity": 0.6,
    "disabled-opacity": 0.38,
    "idle-opacity": 0.04,
    "hover-opacity": 0.04,
    "focus-opacity": 0.12,
    "selected-opacity": 0.08,
    "activated-opacity": 0.12,
    "pressed-opacity": 0.12,
    "dragged-opacity": 0.08,
    "theme-kbd": "#212529",
    "theme-on-kbd": "#FFFFFF",
    "theme-code": "#F5F5F5",
    "theme-on-code": "#000000",
  },
};
// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: "system",
    themes: {
      light: lightTheme,
      dark: darkTheme,
    },
  },
  defaults: {
    VTextField: {
      variant: "outlined",
      density: "comfortable",
      color: "primary",
      rounded: "lg",
      hideDetails: "auto",
    },
    VSelect: {
      variant: "outlined",
      density: "comfortable",
      color: "primary",
      rounded: "lg",
      hideDetails: "auto",
    },
    VSheet: {
      rounded: "lg",
    },
    VCard: {
      rounded: "lg",
    },
    VBtn: {
      rounded: "lg",
      elevation: 0,
      variant: "outlined",
      color: "primary",
    },
    VToolbar: {
      flat: true,
      color: "transparent",
      density: "comfortable",
      rounded: "lg",
      elevation: 1,
    },
  },
});
