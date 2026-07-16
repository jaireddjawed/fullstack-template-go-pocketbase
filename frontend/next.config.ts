import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Produce the self-contained Node server used by the production Docker image.
  output: "standalone",
};

export default nextConfig;
