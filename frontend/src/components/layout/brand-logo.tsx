import { cn } from "@/lib/utils";

export function BrandLogo({ className }: { className?: string }) {
  return (
    <div className={cn("flex items-center gap-3", className)}>
      <div className="relative flex h-11 w-11 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-sm">
        <svg viewBox="0 0 48 48" aria-hidden="true" className="h-7 w-7">
          <path
            d="M12 25c0-7.18 5.82-13 13-13s13 5.82 13 13"
            fill="none"
            stroke="currentColor"
            strokeWidth="3.2"
            strokeLinecap="round"
          />
          <path
            d="M15 24.5v7.5a3.5 3.5 0 0 0 3.5 3.5H21V24.5zM33 24.5v7.5a3.5 3.5 0 0 1-3.5 3.5H27V24.5z"
            fill="currentColor"
          />
          <path
            d="M24 19v9"
            stroke="#f4f1ea"
            strokeWidth="2.8"
            strokeLinecap="round"
          />
        </svg>
      </div>
      <div>
        <p className="text-sm font-semibold text-foreground">Kijani Voice</p>
        <p className="text-xs text-muted-foreground">Call center operations</p>
      </div>
    </div>
  );
}
