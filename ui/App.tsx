import React, { useState } from 'react';
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Upload, MessageSquare, Search, Database } from 'lucide-react'
import { uploadDocument, searchDocuments, generateResponse } from '@/utils/ragUtils';
import { Message, Document } from '@/types';

function App() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [documents, setDocuments] = useState<Document[]>([]);
  const [activeTab, setActiveTab] = useState('chat');
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Document[]>([]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim()) {
      setMessages([...messages, { role: 'user', content: input }]);
      const response = await generateResponse(input, documents);
      setMessages(prev => [...prev, { role: 'assistant', content: response }]);
      setInput('');
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const newDoc = await uploadDocument(file);
      setDocuments([...documents, newDoc]);
    }
  };

  const handleSearch = async () => {
    if (searchQuery.trim()) {
      const results = await searchDocuments(searchQuery, documents);
      setSearchResults(results);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      <header className="bg-white shadow-sm p-4">
        <h1 className="text-xl font-bold">Advanced RAG System</h1>
      </header>
      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="chat">
            <MessageSquare className="w-4 h-4 mr-2" />
            Chat
          </TabsTrigger>
          <TabsTrigger value="documents">
            <Upload className="w-4 h-4 mr-2" />
            Documents
          </TabsTrigger>
          <TabsTrigger value="search">
            <Search className="w-4 h-4 mr-2" />
            Search
          </TabsTrigger>
          <TabsTrigger value="knowledge">
            <Database className="w-4 h-4 mr-2" />
            Knowledge Base
          </TabsTrigger>
        </TabsList>
        <TabsContent value="chat" className="flex-1 overflow-hidden flex flex-col">
          <ScrollArea className="flex-1 p-4">
            {messages.map((message, index) => (
              <div key={index} className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'} mb-4`}>
                <div className={`flex items-start ${message.role === 'user' ? 'flex-row-reverse' : ''}`}>
                  <Avatar className="w-8 h-8">
                    <AvatarImage src={message.role === 'user' ? '/user-avatar.png' : '/ai-avatar.png'} />
                    <AvatarFallback>{message.role === 'user' ? 'U' : 'AI'}</AvatarFallback>
                  </Avatar>
                  <div className={`mx-2 p-3 rounded-lg ${message.role === 'user' ? 'bg-blue-500 text-white' : 'bg-white'}`}>
                    {message.content}
                  </div>
                </div>
              </div>
            ))}
          </ScrollArea>
          <footer className="bg-white p-4">
            <form onSubmit={handleSubmit} className="flex space-x-2">
              <Textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Ask a question..."
                className="flex-1"
              />
              <Button type="submit">Send</Button>
            </form>
          </footer>
        </TabsContent>
        <TabsContent value="documents" className="flex-1 p-4 overflow-auto">
          <Card>
            <CardHeader>
              <CardTitle>Upload Documents</CardTitle>
              <CardDescription>Add documents to the RAG system</CardDescription>
            </CardHeader>
            <CardContent>
              <Input type="file" onChange={handleFileUpload} accept=".txt,.pdf,.doc,.docx" />
              <div className="mt-4">
                <h3 className="font-semibold mb-2">Uploaded Documents:</h3>
                {documents.map((doc) => (
                  <div key={doc.id} className="mb-2">
                    <span className="font-medium">{doc.title}</span>
                    <p className="text-sm text-gray-500">{doc.content.slice(0, 50)}...</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="search" className="flex-1 p-4 overflow-auto">
          <Card>
            <CardHeader>
              <CardTitle>Search Documents</CardTitle>
              <CardDescription>Search through uploaded documents</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex space-x-2 mb-4">
                <Input
                  type="text"
                  placeholder="Enter search query"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
                <Button onClick={handleSearch}>Search</Button>
              </div>
              <div>
                <h3 className="font-semibold mb-2">Search Results:</h3>
                {searchResults.map((doc) => (
                  <div key={doc.id} className="mb-2">
                    <span className="font-medium">{doc.title}</span>
                    <p className="text-sm text-gray-500">{doc.content.slice(0, 100)}...</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="knowledge" className="flex-1 p-4 overflow-auto">
          <Card>
            <CardHeader>
              <CardTitle>Knowledge Base</CardTitle>
              <CardDescription>Manage and view your knowledge base</CardDescription>
            </CardHeader>
            <CardContent>
              <h3 className="font-semibold mb-2">Knowledge Base Statistics:</h3>
              <p>Total Documents: {documents.length}</p>
              <p>Total Words: {documents.reduce((acc, doc) => acc + doc.content.split(' ').length, 0)}</p>
              {/* Add more knowledge base management features here */}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default App;