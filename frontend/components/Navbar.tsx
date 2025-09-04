"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../app/config";
import { useUser } from "../contexts/UserContext";

export default function Navbar() {
  const { userId, setUserId } = useUser();
  const router = useRouter();

  // Vérifie si l'utilisateur est connecté via cookie au montage
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

  const logout = async () => {
    await fetch(`${API_URL}/api/logout`, { method: "POST", credentials: "include" });
    setUserId(null);
    router.push("/login");
  };

  const goToProfile = () => {
    if (userId) router.push(`/profile/${userId}`);
    else router.push("/login");
  };

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <div className="font-bold cursor-pointer" onClick={() => router.push("/feed")}>
        SocialApp
      </div>
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
