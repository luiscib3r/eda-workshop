import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../backend/openapi/api.swagger.json",
  output: "src/api",
  plugins: ["@tanstack/react-query"],
});
