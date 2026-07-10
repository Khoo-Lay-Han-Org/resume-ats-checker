import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  clientPrefix: "PUBLIC_",
  client: {
    PUBLIC_SENTRY_DSN: z.httpUrl(),
  },
  server: {
    SENTRY_SERVER_DSN: z.httpUrl(),
  },
  runtimeEnv: {
    PUBLIC_SENTRY_DSN: process.env.PUBLIC_SENTRY_DSN,
    SENTRY_SERVER_DSN: process.env.SENTRY_SERVER_DSN,
  },
  skipValidation: !!process.env.SKIP_ENV_VALIDATION,
  emptyStringAsUndefined: true,
});
