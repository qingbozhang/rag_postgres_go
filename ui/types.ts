export type Message = {
  role: 'user' | 'assistant';
  content: string;
};

export type Document = {
  id: string;
  title: string;
  content: string;
};