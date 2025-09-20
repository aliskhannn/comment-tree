import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import Comment from '../components/Comment';
import { getComment } from '../api/comments';
import type { Comment as CommentType } from '../types/types';

const CommentThreadPage = () => {
  const { id } = useParams<{ id: string }>();
  const [comment, setComment] = useState<CommentType | null>(null);

  useEffect(() => {
    const fetchComment = async () => {
      if (!id) return;
      try {
        const data = await getComment(id);
        setComment(data);
      } catch (error) {
        console.error('Failed to fetch comment:', error);
      }
    };
    fetchComment();
  }, [id]);

  if (!comment) {
    return <div className="max-w-4xl mx-auto p-4">Loading...</div>;
  }

  return (
    <div className="max-w-4xl mx-auto p-4">
      <Link to="/" className="text-blue-500 hover:text-blue-700 mb-4 block">
        Back to main comments
      </Link>
      <h1 className="text-2xl font-bold mb-6">Comment Thread</h1>
      <Comment comment={comment} />
    </div>
  );
};

export default CommentThreadPage;