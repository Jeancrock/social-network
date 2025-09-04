// app/layout.tsx
import Navbar from "../components/Navbar";
import { UserProvider } from "../contexts/UserContext";
import "./globals.css";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <UserProvider>
      <html lang="en">
        <body>
          <Navbar />
          {children}
        </body>
      </html>
    </UserProvider>
  );
}
