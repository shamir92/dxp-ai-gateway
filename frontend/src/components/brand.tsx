import Link from "next/link";
import { cn } from "@/lib/utils";

export function BrandMark({ className }: { className?: string }) {
  return (
    <span
      className={cn(
        "flex size-8 items-center justify-center rounded-xl bg-flame",
        className,
      )}
    >
      <span className="size-3 rounded-full bg-white" />
    </span>
  );
}

export function Brand({
  className,
  href = "/",
  inverted = false,
}: {
  className?: string;
  href?: string;
  inverted?: boolean;
}) {
  return (
    <Link
      href={href}
      className={cn("flex items-center gap-2.5 font-semibold tracking-tight", className)}
    >
      <BrandMark />
      <span className={cn("text-[17px]", inverted ? "text-white" : "text-ink")}>
        DXP AI Gateway
      </span>
    </Link>
  );
}
