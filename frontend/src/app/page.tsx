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
        Auth uses the PocketBase JS SDK with an httpOnly cookie for SSR.
      </p>
      <p>
        {stats.total} posts — {stats.published} published, {stats.drafts}{" "}
        drafts. (Typed by <code>PostStats</code>, generated from the Go DTO.)
      </p>
    </>
  );
}
