import { router, usePage } from "@inertiajs/react";
import Layout from "@/components/Layout";
import type { PostSummary } from "@shared/types.gen";
import type { SharedProps } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function PostsIndex({ posts }: { posts: PostSummary[] }) {
  const { auth } = usePage<SharedProps>().props;

  return (
    <Layout>
      <div className="space-y-4">
      <h1>Posts</h1>
      <Card className="border-primary/30 bg-primary/5">
        <CardHeader>
          <CardTitle>Protected component example</CardTitle>
          <CardDescription>
            This card is rendered only after the server validates the PocketBase auth cookie.
          </CardDescription>
        </CardHeader>
        <CardContent>
          Signed in as <strong>{auth?.email}</strong>.
        </CardContent>
      </Card>
      {posts.length === 0 && <p className="muted">No posts yet — run `make seed`.</p>}
      {posts.map((post) => (
        <Card key={post.id}>
          <CardHeader><CardTitle>{post.title}</CardTitle><CardDescription>
            /{post.slug} · {post.published ? "published" : "draft"} ·{" "}
            {post.created}
          </CardDescription></CardHeader>
          <CardContent>
          <div dangerouslySetInnerHTML={{ __html: post.content }} />
          {post.isOwner && !post.published && (
            <Button className="mt-4" onClick={() => router.post(`/posts/${post.id}/publish`)}>
              Publish
            </Button>
          )}
          </CardContent>
        </Card>
      ))}
      </div>
    </Layout>
  );
}
