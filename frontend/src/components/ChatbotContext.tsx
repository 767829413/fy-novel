import React, { createContext, useContext, useState, ReactNode } from 'react';

interface ChatbotContextType {
  isInitialized: boolean;
  setIsInitialized: (value: boolean) => void;
  currentModel: string;
  setCurrentModel: (value: string) => void;
  modelList: string[];
  setModelList: (value: string[]) => void;
}

const ChatbotContext = createContext<ChatbotContextType | undefined>(undefined);

export const ChatbotProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isInitialized, setIsInitialized] = useState(false);
  const [currentModel, setCurrentModel] = useState('');
  const [modelList, setModelList] = useState<string[]>([]);

  return (
    <ChatbotContext.Provider value={{
      isInitialized,
      setIsInitialized,
      currentModel,
      setCurrentModel,
      modelList,
      setModelList
    }}>
      {children}
    </ChatbotContext.Provider>
  );
};

export const useChatbotContext = () => {
  const context = useContext(ChatbotContext);
  if (context === undefined) {
    throw new Error('useChatbotContext must be used within a ChatbotProvider');
  }
  return context;
};