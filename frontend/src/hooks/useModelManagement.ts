import { useState, useEffect, useCallback } from 'react';
import { message } from 'antd';
import { useTranslation } from 'react-i18next';
import { GetCurrentUseModel, GetSelectModelList, SetOllamaModel, GetSetOllamaModelProgress } from '../../wailsjs/go/main/App';
import { useModelChange } from '../context/ModelChangeContext';

export const useModelManagement = () => {
  const { t } = useTranslation();
  const { isChangingModel, modelChangeProgress, setIsChangingModel, setModelChangeProgress } = useModelChange();
  const [currentModel, setCurrentModel] = useState('');
  const [modelList, setModelList] = useState<string[]>([]);
  const [selectedModel, setSelectedModel] = useState<string | null>(null);

  const getCurrentModel = useCallback(async () => {
    try {
      const result = await GetCurrentUseModel();
      if (result.Model) {
        setCurrentModel(result.Model);
      } else {
        console.error('No model returned from GetCurrentUseModel');
        message.error(t('chatbot.noModelError'));
      }
    } catch (error) {
      console.error('Error getting current model:', error);
      message.error(t('chatbot.getCurrentModelError'));
    }
  }, [t]);

  const getModelList = useCallback(async () => {
    try {
      const result = await GetSelectModelList();
      if (result.Models) {
        setModelList(result.Models);
      } else {
        console.error('No model list returned from GetSelectModelList');
        message.error(t('chatbot.noModelListError'));
      }
    } catch (error) {
      console.error('Error getting model list:', error);
      message.error(t('chatbot.getModelListError'));
    }
  }, [t]);

  const handleModelChange = useCallback((value: string) => {
    setIsChangingModel(true);
    setModelChangeProgress(0);
    try {
      SetOllamaModel(value);
      checkModelChangeProgress();
    } catch (error) {
      console.error('Error initiating model change:', error);
      message.error(t('chatbot.changeModelError'));
      setIsChangingModel(false);
      setSelectedModel(null);
    }
  }, [t, setIsChangingModel, setModelChangeProgress]);

  const checkModelChangeProgress = useCallback(() => {
    const checkProgress = async () => {
      try {
        const result = await GetSetOllamaModelProgress();
        if (result.Exists) {
          const progressPercentage = Math.min(Math.round((result.Completed / result.Total) * 100), 100);
          setModelChangeProgress(progressPercentage);
          if (progressPercentage >= 100) {
            setIsChangingModel(false);
            getCurrentModel();
            setSelectedModel(null);
            message.success(t('chatbot.modelChangeSuccess'));
            return;
          }
          setTimeout(checkProgress, 2000);
        } else {
          setIsChangingModel(false);
          setSelectedModel(null);
          message.error(t('chatbot.modelChangeNotStarted'));
        }
      } catch (error) {
        console.error('Error checking model change progress:', error);
        message.error(t('chatbot.modelChangeProgressError'));
        setIsChangingModel(false);
        setSelectedModel(null);
      }
    };
    checkProgress();
  }, [getCurrentModel, t, setIsChangingModel, setModelChangeProgress]);

  useEffect(() => {
    getCurrentModel();
    getModelList();
  }, [getCurrentModel, getModelList]);

  return {
    currentModel,
    modelList,
    selectedModel,
    isChangingModel,
    modelChangeProgress,
    setSelectedModel,
    handleModelChange,
  };
};