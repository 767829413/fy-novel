import { useState, useEffect, useCallback } from 'react';
import { message } from 'antd';
import { useTranslation } from 'react-i18next';
import { HasInitOllama, InitOllama, GetInitOllamaProgress } from '../../wailsjs/go/main/App';

export const useOllamaInitialization = () => {
  const { t } = useTranslation();
  const [isInitialized, setIsInitialized] = useState(false);
  const [showInitModal, setShowInitModal] = useState(false);
  const [initProgress, setInitProgress] = useState(0);
  const [isInitializing, setIsInitializing] = useState(false);
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const [hasContainer, setHasContainer] = useState(false);

  const checkInitialization = useCallback(async () => {
    try {
      const result = await HasInitOllama();
      setHasContainer(result.Has);
      setIsInitializing(result.IsInit);
      setIsInitialized(result.Has && !result.IsInit);
    } catch (error) {
      console.error('Error checking Ollama initialization:', error);
      message.error(t('chatbot.checkInitError'));
    }
  }, [t]);

  const handleInitialize = useCallback(async () => {
    setIsInitializing(true);
    setInitProgress(0);
    try {
      InitOllama();
    } catch (error) {
      console.error('Error initializing Ollama:', error);
      message.error(t('chatbot.initError'));
      setIsInitializing(false);
    }
  }, [t]);

  const checkProgress = useCallback(async () => {
    try {
      const result = await GetInitOllamaProgress();
      if (result.Exists) {
        const progressPercentage = Math.min(Math.round((result.Completed / result.Total) * 100), 100);
        setInitProgress(progressPercentage);
        if (progressPercentage >= 100) {
          setIsInitialized(true);
          setShowInitModal(false);
          setIsInitializing(false);
        }
      }
    } catch (error) {
      console.error('Error checking initialization progress:', error);
      message.error(t('chatbot.progressCheckError'));
      setIsInitializing(false);
    }
  }, [t]);

  useEffect(() => {
    const checkStatus = async () => {
      setIsCheckingStatus(true);
      await checkInitialization();
      setIsCheckingStatus(false);
    };

    checkStatus();
  }, [checkInitialization]);

  useEffect(() => {
    if (!isCheckingStatus && !hasContainer) {
      setShowInitModal(true);
    }
  }, [isCheckingStatus, hasContainer]);

  useEffect(() => {
    let intervalId: number | undefined;

    if (isInitializing) {
      intervalId = window.setInterval(checkProgress, 1000);
    }

    return () => {
      if (intervalId !== undefined) {
        window.clearInterval(intervalId);
      }
    };
  }, [isInitializing, checkProgress]);

  return {
    isInitialized,
    showInitModal,
    initProgress,
    isInitializing,
    isCheckingStatus,
    hasContainer,
    handleInitialize,
    setShowInitModal,
  };
};