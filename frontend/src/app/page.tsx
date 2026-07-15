import type { PostStats } from "@shared/types.gen";
import { POCKETBASE_URL } from "@/lib/pocketbase";

export default async function Home() {
  // Custom Go endpoint; response typed by the generated PostStats DTO.
  const res = await fetch(`${POCKETBASE_URL}/api/app/posts/stats`, {
    cache: "no-store",
  });
  const stats = (await res.json()) as PostStats;

  return (
    <>
      <h1>Fullstack Template</h1>
      <p>
        Next.js frontend talking to a PocketBase backend extended with Go.
        Auth is handled by Clerk; the Go backend verifies Clerk session
        tokens and maps them to PocketBase user records.
      </p>
      <p>
        {stats.total} posts — {stats.published} published, {stats.drafts}{" "}
        drafts. (Typed by <code>PostStats</code>, generated from the Go DTO.)
      </p>
    </>
  );
}
