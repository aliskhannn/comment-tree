import axios from "axios";
import type { Comment } from "../types/types";
import { buildCommentTree } from "../utils/buildTree";

const API_URL = "http://localhost:8080/api/comments/";

interface CreateCommentResponse {
  result: Comment;
}

export const createComment = async (
  content: string,
  parentId: string | null
) => {
  console.log("Sending POST request with payload:", { content, parentId }); // Debug log
  const response = await axios.post<CreateCommentResponse>(API_URL, {
    content,
    parent_id: parentId,
  });
  return response.data.result;
};

export const getComments = async (params: {
  parent?: string;
  search?: string;
  limit?: number;
  offset?: number;
}) => {
  const response = await axios.get<Comment[]>(API_URL, { params });
  return buildCommentTree(response.data || []);
};

export const getComment = async (id: string) => {
  const response = await axios.get<Comment>(`${API_URL}${id}`);
  return response.data;
};

export const deleteComment = async (id: string) => {
  await axios.delete(`${API_URL}${id}`);
};
