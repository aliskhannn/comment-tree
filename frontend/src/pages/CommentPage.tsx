import { useEffect, useState } from 'react';
import Comment from '../components/Comment';
import CommentForm from '../components/CommentForm';
import SearchBar from '../components/SearchBar';
import { getComments } from '../api/comments';
import type { Comment as CommentType } from '../types/types';

const CommentPage = () => {
  const [comments, setComments] = useState<CommentType[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(1);
  const limit = 10;

  const fetchComments = async () => {
    try {
      const data = await getComments({
        parent: undefined,
        search: searchTerm || undefined,
        limit,
        offset: (page - 1) * limit,
      });
      setComments(data || []);
    } catch (error) {
      console.error('Failed to fetch comments:', error);
    }
  };

  const handleCommentAdded = (newComment: CommentType) => {
    if (!newComment.parent_id) {
      setComments([newComment, ...comments]);
    } else {
      // Optionally refetch to ensure tree consistency
      fetchComments();
    }
  };

  useEffect(() => {
    fetchComments();
  }, [searchTerm, page]);

  return (
    <div className="max-w-4xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">Comments</h1>
      
      <SearchBar onSearch={setSearchTerm} />
      <CommentForm parentId={null} onSuccess={handleCommentAdded} />

      <div className="mt-6">
        {comments.length > 0 ? (
          comments.map((comment) => (
            <Comment
              key={comment.id}
              comment={comment}
              level={0}
              onCommentAdded={handleCommentAdded}
            />
          ))
        ) : (
          <p className="text-gray-500">No comments found</p>
        )}
      </div>

      <div className="mt-4 flex justify-between">
        <button
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
          className="px-4 py-2 bg-gray-200 rounded-md disabled:bg-gray-100"
        >
          Previous
        </button>
        <button
          onClick={() => setPage((p) => p + 1)}
          disabled={comments.length < limit}
          className="px-4 py-2 bg-gray-200 rounded-md disabled:bg-gray-100"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default CommentPage;