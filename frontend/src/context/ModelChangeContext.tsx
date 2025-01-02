import React, { createContext, useState, useContext, ReactNode } from 'react';

interface ModelChangeContextType {
  isChangingModel: boolean;
  modelChangeProgress: number;
  setIsChangingModel: (isChanging: boolean) => void;
  setModelChangeProgress: (progress: number) => void;
}

const ModelChangeContext = createContext<ModelChangeContextType | undefined>(undefined);

export const ModelChangeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isChangingModel, setIsChangingModel] = useState(false);
  const [modelChangeProgress, setModelChangeProgress] = useState(0);

  return (
    <ModelChangeContext.Provider
      value={{
        isChangingModel,
        modelChangeProgress,
        setIsChangingModel,
        setModelChangeProgress,
      }}
    >
      {children}
    </ModelChangeContext.Provider>
  );
};

export const useModelChange = () => {
  const context = useContext(ModelChangeContext);
  if (context === undefined) {
    throw new Error('useModelChange must be used within a ModelChangeProvider');
  }
  return context;
};

export { ModelChangeContext };