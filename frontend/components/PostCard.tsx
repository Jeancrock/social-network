"use client";

interface PostCardProps {
  post: {
    id: string;
    username: string;
    content: string;
    created: string;
  };
}

export default function PostCard({ post }: PostCardProps) {
  return (
    <div className="border rounded p-4 shadow">
      <div className="font-bold">{post.username}</div>
      <div className="text-sm text-gray-500">{new Date(post.created).toLocaleString()}</div>
      <p className="mt-2">{post.content}</p>
    </div>
  );
}
