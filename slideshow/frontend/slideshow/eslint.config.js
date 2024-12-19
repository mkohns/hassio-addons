import js from "@eslint/js";
import pluginVue from "eslint-plugin-vue";

export default [
  {
    name: "app/files-to-lint",
    files: ["**/*.{js,mjs,jsx,vue}"],
  },

  {
    name: "app/files-to-ignore",
    ignores: ["**/dist/**", "**/dist-ssr/**", "**/coverage/**"],
  },

  js.configs.recommended,
  ...pluginVue.configs["flat/recommended"],

  {
    extends: [
      "eslint:recommended",
      "plugin:vue/vue3-recommended", // Adjust for Vue 2 or other frameworks
      "plugin:prettier/recommended",
    ],
    rules: {
      "vue/multi-word-component-names": "off",
      "vue/html-self-closing": ["off"],
      "vue/max-attributes-per-line": "off",
    },
  },
];
