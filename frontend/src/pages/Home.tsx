import Layout from "@/components/Layout";
import type { PostStats } from "@shared/types.gen";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function Home({ stats }: { stats: PostStats }) {
  return (
    <Layout>
      <Card>
        <CardHeader><CardTitle>Fullstack Template</CardTitle><CardDescription>React, Inertia, PocketBase, and Go.</CardDescription></CardHeader>
        <CardContent className="space-y-4">
        <p>
        Go + PocketBase backend rendering a React frontend over Inertia.js —
        no client-side routing or API glue needed for pages.
      </p>
      <p>
        {stats.total} posts — {stats.published} published, {stats.drafts}{" "}
        drafts. (Typed by <code>PostStats</code>, generated from the Go DTO.)
      </p>
        </CardContent>
      </Card>
    </Layout>
  );
}
