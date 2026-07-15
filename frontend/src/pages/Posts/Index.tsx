import { router } from "@inertiajs/react";
import Layout from "@/components/Layout";
import type { PostSummary } from "@shared/types.gen";

export default function PostsIndex({ posts }: { posts: PostSummary[] }) {
  return (
    <Layout>
      <h1>Posts</h1>
      {posts.length === 0 && <p className="muted">No posts yet — run `make seed`.</p>}
      {posts.map((post) => (
        <article key={post.id}>
          <h2>{post.title}</h2>
          <p className="muted">
            /{post.slug} · {post.published ? "published" : "draft"} ·{" "}
            {post.created}
          </p>
          <div dangerouslySetInnerHTML={{ __html: post.content }} />
          {post.isOwner && !post.published && (
            <button onClick={() => router.post(`/posts/${post.id}/publish`)}>
              Publish
            </button>
          )}
        </article>
      ))}
    </Layout>
  );
}
