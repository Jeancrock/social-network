"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../config";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Payload envoyé:", { username, email, password });
    const res = await fetch(`${API_URL}/api/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });
    if (res.ok) router.push("/login");
    else alert("Registration failed");
  };

  return (
    <form onSubmit={submit} className="max-w-md mx-auto mt-10 space-y-2">
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        className="w-full border p-2 rounded"
      />
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
      <button type="submit" className="w-full bg-green-500 text-white py-2 rounded">
        Register
      </button>
    </form>
  );
}
