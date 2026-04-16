import pino from "pino";

const isDev = import.meta.env.DEV;

export const logger = pino({
  level: isDev ? "debug" : "warn",
  browser: {
    // log to console.* with nice objects in devtools
    asObject: true,

    // optional: write mapping to console methods
    write:  {
      info (o) {
          console.info(o);
      },
      debug (o) {
          console.debug(o);
      },
      warn (o) {
          console.warn(o);
      },
      error (o) {
          console.error(o);
      },
    },
  },
  base: {
    app: "skemr-frontend",
    version: "0.0.1",
  },
});
