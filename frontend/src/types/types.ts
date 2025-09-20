export interface Comment {
  id: string;
  parent_id: string | null;
  content: string;
  created_at: string;
  updated_at: string;
  children?: Comment[]; // This will be added by our tree builder
}