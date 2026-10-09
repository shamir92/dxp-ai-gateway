import Link from "next/link";
import { cn } from "@/lib/utils";

export function BrandMark({
  className,
  inverted = false,
}: {
  className?: string;
  inverted?: boolean;
}) {
  return (
    <img
      src={inverted ? "/dxp-logo-light.svg" : "/dxp-logo.svg"}
      alt="DXP."
      width={72}
      height={29}
      className={cn("h-7 w-auto", className)}
    />
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
      <BrandMark inverted={inverted} className="h-8 w-auto" />
      <span className={cn("text-[17px]", inverted ? "text-white" : "text-ink")}>
        AI Gateway
      </span>
    </Link>
  );
}
