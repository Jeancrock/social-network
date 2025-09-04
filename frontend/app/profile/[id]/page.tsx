"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { API_URL } from "../../config";
import { useUser } from "../../../contexts/UserContext";

interface Post {
  id: string;
  userId: string;
  content: string;
  created: string;
}

interface User {
  id: string;
  username: string;
  email: string;
}

interface ProfileResponse {
  user: User;
  posts: Post[];
  followers: User[];
  following: User[];
  groups: string[];
}

export default function ProfilePage() {
  const params = useParams();
  const id = params.id as string;
  const router = useRouter();
  const { userId } = useUser();

  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [isFollowing, setIsFollowing] = useState<boolean | null>(null);

  // Vérifie si connecté
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(`${API_URL}/api/profile?id=${id}`, {
          credentials: "include",
        });
        if (!res.ok) {
          router.push("/login"); // redirection si pas connecté
          return;
        }
        const data: ProfileResponse = await res.json();
        setProfile(data);

        // Détermine si l'utilisateur suit ce profil
        if (userId && userId !== data.user.id) {
          setIsFollowing(data.followers.some(f => f.id === userId));
        }
      } catch {
        router.push("/login");
      } finally {
        setLoading(false);
      }
    };
    checkAuth();
  }, [id, router, userId]);

  const toggleFollow = async () => {
    if (!userId || userId === profile?.user.id) return;
    try {
      const method = isFollowing ? "DELETE" : "POST";
      const res = await fetch(`${API_URL}/api/follow?userId=${id}`, {
        method,
        credentials: "include",
      });
      if (res.ok && profile) {
        setIsFollowing(!isFollowing);
        // Met à jour la liste des followers localement
        setProfile({
          ...profile,
          followers: isFollowing
            ? profile.followers.filter(f => f.id !== userId)
            : [...profile.followers, { id: userId, username: "You", email: "" }],
        });
      }
    } catch (err) {
      console.error("Failed to toggle follow", err);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (!profile) return <div>Profile not found</div>;

  return (
    <div className="p-4">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">{profile.user.username}</h1>
        {userId === profile.user.id ? (
          <button
            className="bg-gray-500 text-white px-3 py-1 rounded"
            onClick={() => router.push("/settings")}
          >
            Settings
          </button>
        ) : (
          isFollowing !== null && (
            <button
              className={`px-3 py-1 rounded text-white ${
                isFollowing ? "bg-red-500" : "bg-blue-500"
              }`}
              onClick={toggleFollow}
            >
              {isFollowing ? "Unfollow" : "Follow"}
            </button>
          )
        )}
      </div>

      <h2 className="mt-4 font-semibold">Posts</h2>
      {profile.posts.length === 0 ? (
        <p>No posts</p>
      ) : (
        <ul>
          {profile.posts.map((p) => (
            <li key={p.id} className="mb-2 border p-2 rounded">
              <p>{p.content}</p>
              <small>{new Date(p.created).toLocaleString()}</small>
            </li>
          ))}
        </ul>
      )}

      <h2 className="mt-4 font-semibold">Followers ({profile.followers.length})</h2>
      <ul>
        {profile.followers.map((f) => (
          <li key={f.id}>{f.username}</li>
        ))}
      </ul>

      <h2 className="mt-4 font-semibold">Following ({profile.following.length})</h2>
      <ul>
        {profile.following.map((f) => (
          <li key={f.id}>{f.username}</li>
        ))}
      </ul>

      <h2 className="mt-4 font-semibold">Groups ({profile.groups.length})</h2>
      <ul>
        {profile.groups.map((g, i) => (
          <li key={i}>{g}</li>
        ))}
      </ul>
    </div>
  );
}
