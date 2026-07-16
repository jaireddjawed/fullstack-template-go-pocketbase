import { Collections } from "@shared/pocketbase.gen";
import { createServerClient, currentUser } from "@/lib/server-auth";
import PublishButton from "@/components/PublishButton";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

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
    <div className="space-y-4">
      <h1>Posts</h1>
      <Card className="border-primary/30 bg-primary/5">
        <CardHeader>
          <CardTitle>Protected component example</CardTitle>
          <CardDescription>
            This card is rendered only after Clerk validates the session and the Go backend provisions the viewer.
          </CardDescription>
        </CardHeader>
        <CardContent>
          Signed in as <strong>{user?.email}</strong>.
        </CardContent>
      </Card>
      {posts.length === 0 && (
        <p className="muted">No posts yet — run `make seed`.</p>
      )}
      {posts.map((post) => (
        <Card key={post.id}>
          <CardHeader>
            <CardTitle>{post.title}</CardTitle>
            <CardDescription>
            /{post.slug} · {post.published ? "published" : "draft"} ·{" "}
            {post.created}
            </CardDescription>
          </CardHeader>
          <CardContent>
          <div dangerouslySetInnerHTML={{ __html: post.content }} />
          {post.owner === userID && !post.published && (
            <div className="mt-4"><PublishButton id={post.id} /></div>
          )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
