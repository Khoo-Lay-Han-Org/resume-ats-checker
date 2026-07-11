import node from "@astrojs/node";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, envField } from "astro/config";
import sentry from "@sentry/astro";
import svelte from "@astrojs/svelte";

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
      PUBLIC_SENTRY_DSN: envField.string({
        access: "public",
        context: "client",
        optional: true,
      }),
      SENTRY_SERVER_DSN: envField.string({
        access: "public",
        context: "server",
        optional: true,
      })
    },
  },

  vite: {
    plugins: [tailwindcss()],
  },

  integrations: [
    sentry({
      project: process.env.SENTRY_PROJECT_NAME,
      org: process.env.SENTRY_PROJECT_ORG,
      authToken: process.env.SENTRY_AUTH_TOKEN,
    }), 
    svelte(),
  ],
});
