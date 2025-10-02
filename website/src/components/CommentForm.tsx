import { useState } from 'react';
import { createComment } from '../api/comments';
import type { Comment } from '../types/types';

interface CommentFormProps {
  parentId?: string | null; // parentId can be string or null
  onSuccess?: (newComment: Comment) => void;
}

const CommentForm: React.FC<CommentFormProps> = ({ parentId, onSuccess }) => {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;

    setIsSubmitting(true);
    try {
      console.log('Submitting comment with parentId:', parentId); // Debug log
      const newComment = await createComment(content, parentId ?? null);
      setContent('');
      onSuccess?.(newComment);
    } catch (error) {
      console.error('Failed to create comment:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mt-4">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        className="w-full p-2 border rounded-md"
        rows={3}
        placeholder="Write your comment..."
        disabled={isSubmitting}
      />
      <button
        type="submit"
        className="mt-2 bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 disabled:bg-blue-300"
        disabled={isSubmitting}
      >
        {isSubmitting ? 'Submitting...' : 'Post Comment'}
      </button>
    </form>
  );
};

export default CommentForm;