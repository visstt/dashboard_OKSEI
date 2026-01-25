import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": "/src",
      "@lib": "/src/lib",
      "@components": "/src/components",
      "@pages": "/src/pages",
    },
  },
});
