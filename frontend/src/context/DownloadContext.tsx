import React, { createContext, useState, useContext } from 'react';

interface DownloadContextType {
  isDownloading: boolean;
  setIsDownloading: (isDownloading: boolean) => void;
}

const DownloadContext = createContext<DownloadContextType | undefined>(undefined);

export const DownloadProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isDownloading, setIsDownloading] = useState(false);

  return (
    <DownloadContext.Provider value={{ isDownloading, setIsDownloading }}>
      {children}
    </DownloadContext.Provider>
  );
};

export const useDownload = () => {
  const context = useContext(DownloadContext);
  if (context === undefined) {
    throw new Error('useDownload must be used within a DownloadProvider');
  }
  return context;
};