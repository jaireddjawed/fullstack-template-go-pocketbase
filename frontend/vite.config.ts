import { defineConfig, type Plugin } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const dirname = path.dirname(fileURLToPath(import.meta.url));
const hotFile = path.join(dirname, ".hot");

// Laravel-style "hot file": while the dev server runs, frontend/.hot holds
// its URL. The Go backend (gonertia) sees the file and points asset tags at
// the dev server; when it's gone, it falls back to dist/manifest.json.
function hotFilePlugin(): Plugin {
  return {
    name: "go-hot-file",
    configureServer(server) {
      server.httpServer?.once("listening", () => {
        const address = server.httpServer?.address();
        if (address && typeof address === "object") {
          fs.writeFileSync(hotFile, `http://localhost:${address.port}`);
        }
      });
      const remove = () => fs.existsSync(hotFile) && fs.rmSync(hotFile);
      process.on("exit", remove);
      process.on("SIGINT", () => process.exit());
      process.on("SIGTERM", () => process.exit());
    },
  };
}

export default defineConfig({
  plugins: [react(), tailwindcss(), hotFilePlugin()],
  resolve: {
    alias: {
      "@": path.resolve(dirname, "src"),
      "@shared": path.resolve(dirname, "../shared"),
    },
  },
  server: {
    cors: true, // the page is served by the Go backend on another port
  },
  build: {
    manifest: true,
    outDir: "dist",
    rollupOptions: {
      input: "src/main.tsx",
    },
  },
});
