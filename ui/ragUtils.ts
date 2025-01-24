import { Document } from '@/types';

// Simulated document upload function
export const uploadDocument = async (file: File): Promise<Document> => {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      const content = event.target?.result as string;
      resolve({
        id: Date.now().toString(),
        title: file.name,
        content: content,
      });
    };
    reader.readAsText(file);
  });
};

// Simulated document search function
export const searchDocuments = async (query: string, documents: Document[]): Promise<Document[]> => {
  // Simple search implementation. In a real system, this would use more sophisticated search algorithms.
  return documents.filter(doc =>
    doc.content.toLowerCase().includes(query.toLowerCase()) ||
    doc.title.toLowerCase().includes(query.toLowerCase())
  );
};

// Simulated response generation function
export const generateResponse = async (query: string, documents: Document[]): Promise<string> => {
  // Simple response generation. In a real system, this would use a language model and retrieved documents.
  const relevantDocs = await searchDocuments(query, documents);

  if (relevantDocs.length > 0) {
    return `Based on the documents I've found, here's a response to "${query}":
    ${relevantDocs[0].content.slice(0, 200)}...
    This information comes from the document titled "${relevantDocs[0].title}".`;
  } else {
    return `I'm sorry, I couldn't find any relevant information to answer "${query}".
    Can you please provide more context or rephrase your question?`;
  }
};