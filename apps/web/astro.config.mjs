import node from "@astrojs/node";
// @ts-check
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, envField } from "astro/config";

import sentry from "@sentry/astro";

// https://astro.build/config
export default defineConfig({
  output: "server",
  adapter: node({ mode: "standalone" }),

  env: {
    schema: {
      PUBLIC_SERVER_URL: envField.string({
        access: "public",
        context: "client",
        default: "http://localhost:3000",
      }),
      SENTRY_PROJECT_NAME: envField.string({
        access: "public",
        context: "server",
      }),
      SENTRY_PROJECT_ORG: envField.string({
        access: "public",
        context: "server",
      }),
    },
  },

  vite: {
    plugins: [tailwindcss()],
  },

  integrations: [sentry({
    project: process.env.SENTRY_PROJECT_NAME,
    org: process.env.SENTRY_PROJECT_ORG,
    authToken: process.env.SENTRY_AUTH_TOKEN,
  })],
});