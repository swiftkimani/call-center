import type { ReactNode } from "react";
import type { Metadata } from "next";
import { Manrope, JetBrains_Mono } from "next/font/google";

import "./globals.css";
import { QueryProvider } from "@/components/providers/query-provider";

const manrope = Manrope({
  subsets: ["latin"],
  variable: "--font-sans",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "Call Center Command",
  description: "Next.js 16 dashboard for the Go call center platform.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${manrope.variable} ${jetbrainsMono.variable}`}>
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
