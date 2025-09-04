"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../config";

interface Post {
  id: string;
  userId: string;
  username: string;
  content: string;
  created: string;
}

export default function FeedPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [newPost, setNewPost] = useState("");
  const [loading, setLoading] = useState(true); // état pour vérifier auth
  const router = useRouter();

  // Vérifie si connecté avant d'afficher quoi que ce soit
  const checkAuth = async () => {
    try {
      const res = await fetch(`${API_URL}/api/profile`, { credentials: "include" });
      if (!res.ok) {
        router.push("/login"); // redirection immédiate
        return false;
      }
      return true;
    } catch {
      router.push("/login");
      return false;
    }
  };

  const fetchPosts = async () => {
    try {
      const res = await fetch(`${API_URL}/api/posts`, { credentials: "include" });
      if (!res.ok) throw new Error("Failed to fetch posts");
      const data: Post[] = await res.json();
      setPosts(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error("Fetch posts error:", err);
      setPosts([]);
    }
  };

  const submitPost = async () => {
    if (!newPost.trim()) return;
    try {
      const res = await fetch(`${API_URL}/api/posts`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: newPost }),
        credentials: "include",
      });
      if (res.ok) {
        setNewPost("");
        fetchPosts();
      } else {
        const text = await res.text();
        alert(`Failed to post: ${text}`);
      }
    } catch (err) {
      console.error("Submit post error:", err);
      alert("Failed to post: network error");
    }
  };

  useEffect(() => {
    const init = async () => {
      const isAuth = await checkAuth();
      if (isAuth) {
        await fetchPosts();
        setLoading(false); // prêt à afficher le feed
      }
    };
    init();
  }, []);

  const goToProfile = (userId: string) => {
    router.push(`/profile/${userId}`);
  };

  if (loading) {
    // Rendu bloqué tant que auth non vérifiée
    return <div className="text-center mt-20">Loading...</div>;
  }

  return (
    <div className="max-w-xl mx-auto mt-10 space-y-4">
      <div className="flex space-x-2">
        <input
          type="text"
          placeholder="What's on your mind?"
          value={newPost}
          onChange={(e) => setNewPost(e.target.value)}
          className="flex-1 border p-2 rounded"
        />
        <button
          onClick={submitPost}
          className="bg-green-500 text-white py-2 px-4 rounded"
        >
          Post
        </button>
      </div>

      <div className="space-y-4">
        {(posts || []).map((p) => (
          <div key={p.id} className="border p-2 rounded">
            <strong
              className="cursor-pointer text-blue-500 hover:underline"
              onClick={() => goToProfile(p.userId)}
            >
              {p.username || "Unknown"}
            </strong>
            <p>{p.content}</p>
            <small>
              {p.created ? new Date(p.created).toLocaleString() : "Unknown date"}
            </small>
          </div>
        ))}
      </div>
    </div>
  );
}
