import { Collections } from "@shared/pocketbase.gen";
import { createServerClient, currentUser } from "@/lib/server-auth";
import PublishButton from "@/components/PublishButton";

export const dynamic = "force-dynamic";

export default async function PostsPage() {
  const pb = await createServerClient();
  // The backend maps the Clerk session to a PocketBase user record;
  // /api/app/me tells us which one, so we can compare against post.owner.
  const user = await currentUser();
  const userID = user?.id;

  // Built-in PocketBase CRUD API; the collection's API rules already limit
  // results to published posts + the viewer's own. Fully typed by
  // shared/pocketbase.gen.ts.
  const posts = await pb.collection(Collections.Posts).getFullList({
    sort: "-created",
  });

  return (
    <>
      <h1>Posts</h1>
      {posts.length === 0 && (
        <p className="muted">No posts yet — run `make seed`.</p>
      )}
      {posts.map((post) => (
        <article key={post.id}>
          <h2>{post.title}</h2>
          <p className="muted">
            /{post.slug} · {post.published ? "published" : "draft"} ·{" "}
            {post.created}
          </p>
          <div dangerouslySetInnerHTML={{ __html: post.content }} />
          {post.owner === userID && !post.published && (
            <PublishButton id={post.id} />
          )}
        </article>
      ))}
    </>
  );
}
