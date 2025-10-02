import { useState } from "react";
import { Link } from "react-router-dom";
import { deleteComment } from "../api/comments";
import type { Comment as CommentType } from "../types/types";
import CommentForm from "./CommentForm";

interface CommentProps {
  comment: CommentType;
  level?: number;
  onCommentAdded?: (newComment: CommentType) => void;
}

const Comment: React.FC<CommentProps> = ({
  comment,
  level = 0,
  onCommentAdded,
}) => {
  const [showReplyForm, setShowReplyForm] = useState(false);
  const [showChildren, setShowChildren] = useState(false);
  const [isDeleted, setIsDeleted] = useState(false);
  const [localComment, setLocalComment] = useState<CommentType>({
    ...comment,
    children: comment.children || [],
  });

  const handleDelete = async () => {
    try {
      await deleteComment(comment.id);
      setIsDeleted(true);
    } catch (error) {
      console.error("Failed to delete comment:", error);
    }
  };

  const handleCommentAdded = (newComment: CommentType) => {
    if (newComment.parent_id === localComment.id) {
      setLocalComment({
        ...localComment,
        children: [newComment, ...(localComment.children || [])],
      });
    }
    setShowReplyForm(false);
    onCommentAdded?.(newComment);
  };

  if (isDeleted) return null;

  const hasChildren = localComment.children && localComment.children.length > 0;
  const hasManyChildren = hasChildren && localComment.children!.length > 4;

  return (
    <div className={`ml-${level * 4} p-4 border-l-2 border-gray-200`}>
      <div className="flex justify-between items-start mb-2">
        <div className="flex-1">
          <p className="text-gray-800">{localComment.content}</p>
          <p className="text-sm text-gray-500">
            {new Date(localComment.created_at).toLocaleString()}
          </p>
        </div>
        <button
          onClick={handleDelete}
          className="text-red-500 hover:text-red-700 ml-2"
        >
          Delete
        </button>
      </div>

      <div className="flex gap-2 flex-wrap">
        <button
          onClick={() => setShowReplyForm(!showReplyForm)}
          className="text-blue-500 hover:text-blue-700 text-sm"
        >
          {showReplyForm ? "Cancel" : "Reply"}
        </button>

        {hasChildren && (
          <button
            onClick={() => setShowChildren(!showChildren)}
            className="text-green-500 hover:text-green-700 text-sm"
          >
            {showChildren ? "Hide" : "Show"} {localComment.children!.length}{" "}
            {localComment.children!.length === 1 ? "reply" : "replies"}
          </button>
        )}
      </div>

      {showReplyForm && (
        <CommentForm
          parentId={localComment.id}
          onSuccess={handleCommentAdded}
        />
      )}

      {hasChildren && showChildren && (
        <div className="mt-4">
          {hasManyChildren ? (
            <Link
              to={`/comment/${localComment.id}`}
              className="text-blue-500 hover:text-blue-700 block mb-2"
            >
              View all {localComment.children!.length} replies in separate
              page...
            </Link>
          ) : (
            localComment.children!.map((child) => (
              <Comment
                key={child.id}
                comment={child}
                level={level + 1}
                onCommentAdded={onCommentAdded}
              />
            ))
          )}
        </div>
      )}
    </div>
  );
};

export default Comment;
