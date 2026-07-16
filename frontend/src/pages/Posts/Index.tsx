import { router } from "@inertiajs/react";
import Layout from "@/components/Layout";
import type { PostSummary } from "@shared/types.gen";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function PostsIndex({ posts }: { posts: PostSummary[] }) {
  return (
    <Layout>
      <div className="space-y-4">
      <h1>Posts</h1>
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
