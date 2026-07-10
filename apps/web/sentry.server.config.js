import * as Sentry from "@sentry/astro";

Sentry.init({
  dsn: import.meta.env.SENTRY_SERVER_DSN,
  dataCollection: {
    // To disable sending user data and HTTP bodies, uncomment the lines below. For more info visit:
    // https://docs.sentry.io/platforms/javascript/guides/astro/configuration/options/#dataCollection
    // userInfo: false,
    // httpBodies: [],
  },
});