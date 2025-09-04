"use client";

import { useState } from "react";
import { API_URL } from "../app/config";

type Props = { onPost: () => void };

export default function NewPostForm({ onPost }: Props) {
  const [content, setContent] = useState("");

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();

    // ✅ Ajout de credentials: "include" pour envoyer le cookie de session
    const res = await fetch(`${API_URL}/api/posts`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ content }),
      credentials: "include",
    });

    if (res.ok) {
      setContent("");
      onPost();
    } else if (res.status === 401) {
      alert("You must be logged in to post"); // message plus précis
    } else {
      alert("Failed to post");
    }
  };

  return (
    <form onSubmit={submit} className="mt-4">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Write something..."
        className="w-full border p-2 rounded"
      />
      <button type="submit" className="mt-2 bg-blue-500 text-white px-4 py-2 rounded">
        Post
      </button>
    </form>
  );
}
