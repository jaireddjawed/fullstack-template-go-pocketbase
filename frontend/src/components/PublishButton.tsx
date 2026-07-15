"use client";

import { useTransition } from "react";
import { publishPost } from "@/lib/actions";

export default function PublishButton({ id }: { id: string }) {
  const [pending, startTransition] = useTransition();

  return (
    <button
      disabled={pending}
      onClick={() => startTransition(() => publishPost(id).then(() => {}))}
    >
      {pending ? "Publishing…" : "Publish"}
    </button>
  );
}
