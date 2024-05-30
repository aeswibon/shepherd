import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "./providers";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full">
      <body className={inter.className}>
        <div className="flex flex-col-reverse md:flex-row md:h-screen relative overflow-x-hidden md:w-screen">
          <Providers>{children}</Providers>
        </div>
      </body>
    </html>
  );
}
