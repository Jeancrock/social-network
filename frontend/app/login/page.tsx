"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../config";
import { useUser } from "../../contexts/UserContext";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();
  const { setUserId } = useUser();

  // Vérifie si déjà connecté au chargement
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(`${API_URL}/api/profile`, { credentials: "include" });
        if (res.ok) {
          const data = await res.json();
          setUserId(data.user.id); // met à jour le contexte
          router.push("/feed"); // redirection automatique
        }
      } catch (err) {
        // pas connecté, on ne fait rien
      }
    };
    checkAuth();
  }, [router, setUserId]);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await fetch(`${API_URL}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    if (res.ok) {
      // Récupérer l'ID du profil connecté
      const profileRes = await fetch(`${API_URL}/api/profile`, { credentials: "include" });
      if (profileRes.ok) {
        const data = await profileRes.json();
        setUserId(data.user.id);
      }
      router.push("/feed");
    } else {
      alert("Login failed");
    }
  };

  return (
    <form onSubmit={submit} className="max-w-md mx-auto mt-10 space-y-2">
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="w-full border p-2 rounded"
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        className="w-full border p-2 rounded"
      />
      <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded">
        Login
      </button>
    </form>
  );
}
