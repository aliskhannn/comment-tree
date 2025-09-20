import type { Comment } from '../types/types';

export const buildCommentTree = (comments: Comment[]): Comment[] => {
  const commentMap = new Map<string, Comment>();
  
  // Add all comments to map
  comments.forEach(comment => {
    comment.children = [];
    commentMap.set(comment.id, { ...comment });
  });

  const roots: Comment[] = [];
  
  // Build tree structure
  comments.forEach(comment => {
    if (comment.parent_id === null) {
      roots.push(commentMap.get(comment.id)!);
    } else {
      const parent = commentMap.get(comment.parent_id);
      if (parent) {
        parent.children!.push(commentMap.get(comment.id)!);
      }
    }
  });

  // Sort children by creation date (newest first)
  roots.forEach(root => sortChildren(root));

  return roots;
};

const sortChildren = (comment: Comment) => {
  if (comment.children && comment.children.length > 0) {
    comment.children.sort((a, b) => 
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    );
    comment.children.forEach(sortChildren);
  }
};