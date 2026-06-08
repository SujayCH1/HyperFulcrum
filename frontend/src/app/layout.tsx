import type { Metadata } from "next";
import { Toaster } from "sonner";
// @ts-ignore: CSS side-effect import
import "./globals.css";

export const metadata: Metadata = {
  title: "HyperFulcrum – SQL Sharding Orchestration",
  description: "Distribute your database across infinite nodes",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="min-h-screen bg-[#050505] text-zinc-400 font-sans selection:bg-zinc-800 selection:text-zinc-100">
        {children}
        <Toaster
          theme="dark"
          position="bottom-right"
          toastOptions={{ style: { background: "#111", border: "1px solid #27272a", color: "#e4e4e7" } }}
        />
      </body>
    </html>
  );
}
