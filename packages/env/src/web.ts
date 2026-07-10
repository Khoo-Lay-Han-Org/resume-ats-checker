import { createEnv } from "@t3-oss/env-core";

export const env = createEnv({
  clientPrefix: "PUBLIC_",
  client: {
  },
  server: {
  },
  runtimeEnv: {
  },
  skipValidation: !!process.env.SKIP_ENV_VALIDATION,
  emptyStringAsUndefined: true,
});
