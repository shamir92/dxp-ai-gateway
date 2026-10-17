import Link from "next/link";
import { cn } from "@/lib/utils";

function GatewayMark() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 16 16"
      fill="none"
      className="size-[0.85em] shrink-0 text-[#ff5a2a]"
    >
      <path
        d="M10.5 3.5h3v9h-3M2 8h9m-3-3 3 3-3 3"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export function BrandMark({
  className,
  inverted = false,
}: {
  className?: string;
  inverted?: boolean;
}) {
  return (
    <span
      role="img"
      aria-label="DXP AI Gateway"
      className={cn(
        "inline-flex items-center gap-1.5 whitespace-nowrap text-base font-semibold leading-none tracking-tight",
        inverted ? "text-[#f3efe6]" : "text-ink",
        className,
      )}
    >
      <span>DXP</span>
      <GatewayMark />
      <span className="site-brand-name">AI Gateway</span>
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
    <Link href={href} className={cn("inline-flex items-center", className)}>
      <BrandMark inverted={inverted} className="text-[18px]" />
    </Link>
  );
}
