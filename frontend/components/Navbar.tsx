"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { API_URL } from "../app/config";
import { useUser } from "../contexts/UserContext";

export default function Navbar() {
  const { userId, setUserId } = useUser();
  const router = useRouter();
  const pathname = usePathname();

  const [query, setQuery] = useState("");
  const [results, setResults] = useState<{ users: any[]; groups: any[] }>({ users: [], groups: [] });

  // Vérifie si connecté au montage
  useEffect(() => {
    fetch(`${API_URL}/api/profile`, { credentials: "include" })
      .then(async (res) => {
        if (res.ok) {
          const data = await res.json();
          setUserId(data.user.id);
        } else {
          setUserId(null);
        }
      })
      .catch(() => setUserId(null));
  }, [setUserId]);

  // Vide la searchbar à chaque changement de page
  useEffect(() => {
    setQuery("");
    setResults({ users: [], groups: [] });
  }, [pathname]);

  const logout = async () => {
    await fetch(`${API_URL}/api/logout`, { method: "POST", credentials: "include" });
    setUserId(null);
    router.push("/login");
  };

  const goToProfile = () => {
    if (userId) router.push(`/profile/${userId}`);
    else router.push("/login");
  };

  const handleSearch = async (q: string) => {
    setQuery(q);
    if (q.trim().length < 2) {
      setResults({ users: [], groups: [] });
      return;
    }
    const res = await fetch(`${API_URL}/api/search?query=${encodeURIComponent(q)}`, {
      credentials: "include",
    });
    if (res.ok) {
      setResults(await res.json());
    }
  };

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center relative">
      <div className="font-bold cursor-pointer" onClick={() => router.push("/feed")}>
        SocialApp
      </div>

      {userId && (
        <div className="relative w-64">
          <input
            type="text"
            placeholder="Search users or groups..."
            value={query}
            onChange={(e) => handleSearch(e.target.value)}
            className="w-full px-3 py-1 rounded text-black"
          />
          {(results.users.length > 0 || results.groups.length > 0) && (
            <div className="absolute bg-white text-black mt-1 w-full rounded shadow-lg z-50">
              {results.users.map((u) => (
                <div
                  key={u.id}
                  className="p-2 hover:bg-gray-200 cursor-pointer"
                  onClick={() => router.push(`/profile/${u.id}`)}
                >
                  👤 {u.username}
                </div>
              ))}
              {results.groups.map((g) => (
                <div
                  key={g.id}
                  className="p-2 hover:bg-gray-200 cursor-pointer"
                  onClick={() => router.push(`/group/${g.id}`)}
                >
                  👥 {g.name}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      <div className="space-x-4">
        {userId ? (
          <>
            <button onClick={goToProfile} className="bg-gray-500 px-3 py-1 rounded">
              Profile
            </button>
            <button onClick={logout} className="bg-red-500 px-3 py-1 rounded">
              Logout
            </button>
          </>
        ) : (
          <>
            <button onClick={() => router.push("/login")} className="bg-blue-500 px-3 py-1 rounded">
              Login
            </button>
            <button onClick={() => router.push("/register")} className="bg-green-500 px-3 py-1 rounded">
              Register
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
