"use client";

import { useTransition } from "react";
import { publishPost } from "@/lib/actions";
import { Button } from "@/components/ui/button";

export default function PublishButton({ id }: { id: string }) {
  const [pending, startTransition] = useTransition();

  return (
    <Button
      disabled={pending}
      onClick={() => startTransition(() => publishPost(id).then(() => {}))}
    >
      {pending ? "Publishing…" : "Publish"}
    </Button>
  );
}
