import Layout from "@/components/Layout";
import type { PostStats } from "@shared/types.gen";

export default function Home({ stats }: { stats: PostStats }) {
  return (
    <Layout>
      <h1>Fullstack Template</h1>
      <p>
        Go + PocketBase backend rendering a React frontend over Inertia.js —
        no client-side routing or API glue needed for pages.
      </p>
      <p>
        {stats.total} posts — {stats.published} published, {stats.drafts}{" "}
        drafts. (Typed by <code>PostStats</code>, generated from the Go DTO.)
      </p>
    </Layout>
  );
}
