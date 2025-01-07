import { useState, useCallback } from 'react';
import { message } from 'antd';
import { useTranslation } from 'react-i18next';
import { HasInitOllama, GetInitOllamaProgress, InitOllama, GetCurrentUseModel, GetSelectModelList, SetOllamaModel, GetSetOllamaModelProgress, InitSetOllamaModelTask } from '../../wailsjs/go/main/App';
import { useModelChange } from '../context/ModelChangeContext';

export const useOllama = () => {
    const { t } = useTranslation();
    // 初始化状态
    const [hasContainer, setHasContainer] = useState(false);
    const [showInitModel, setShowInitModel] = useState(false);
    const [initProgress, setInitProgress] = useState(0);
    const [isInitializing, setIsInitializing] = useState(false);
    const [isCheckingStatus, setIsCheckingStatus] = useState(false);
    // 切换model
    const { isChangingModel, modelChangeProgress, setIsChangingModel, setModelChangeProgress } = useModelChange();
    const [currentModel, setCurrentModel] = useState('');
    const [modelList, setModelList] = useState<string[]>([]);
    const [selectedModel, setSelectedModel] = useState<string | null>(null);

    // 初始化
    const handleInitialize = useCallback(async () => {
        if (!isInitializing) {
            setIsInitializing(true);
            setInitProgress(0);
            InitOllama();
        }
    }, [t]);

    // 检查初始化进度
    const checkProgress = useCallback(async () => {
        const result = await GetInitOllamaProgress();
        console.info("Checking init progress: ", result)
        if (result.Exists) {
            const progressPercentage = Math.min(Math.round((result.Completed / result.Total) * 100), 100);
            setInitProgress(progressPercentage);
            if (progressPercentage >= 100) {
                setShowInitModel(false);
                setIsInitializing(false);
                setIsCheckingStatus(false);
            }
        }
        checkCurrentStatus()
    }, [t]);

    // 检查当前所有状态
    const checkCurrentStatus = useCallback(async () => {
        const result = await HasInitOllama();
        console.info(result)
        console.info({
            "initProgress": initProgress, "isInitializing": isInitializing, "isCheckingStatus": isCheckingStatus, "isChangingModel": isChangingModel, "modelChangeProgress": modelChangeProgress
        })
        // 没有容器
        setShowInitModel(!result.Has);
        setHasContainer(result.Has)
        if (result.Has) {
            // 正在初始化
            setIsInitializing(result.IsInit);
            // 需要检查状态
            setIsCheckingStatus(result.IsInit || result.IsSetModel);
            // 正在更改model
            setIsChangingModel(result.IsSetModel)
        }
    }, [t]);

    // 获取当前模型
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

    // 获取模型列表
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

    // 修改模型
    const handleModelChange = useCallback(async (value: string) => {
        if (!isChangingModel) {
            const result = await InitSetOllamaModelTask()
            if (result.ErrorMsg.length > 0) {
                message.error(t('chatbot.errorMessage'));
                console.error(result.ErrorMsg)
                return;
            }
            setIsChangingModel(true);
            setIsCheckingStatus(true)
            setModelChangeProgress(0);
            SetOllamaModel(value);
            checkModelChangeProgress();
        }
    }, [t, setIsChangingModel, setModelChangeProgress]);


    // 检查模型切换进度
    const checkModelChangeProgress = useCallback(async () => {
        const result = await GetSetOllamaModelProgress();
        console.info("Checking model change progress: ", result)
        if (result.Exists) {
            const progressPercentage = Math.min(Math.round((result.Completed / result.Total) * 100), 100);
            setModelChangeProgress(progressPercentage);
            if (progressPercentage >= 100) {
                setIsChangingModel(false);
                setIsCheckingStatus(false)
                getCurrentModel();
                setSelectedModel(null);
                message.success(t('chatbot.modelChangeSuccess'));
                return;
            }
        }
        checkCurrentStatus()
    }, [t]);


    return {
        // 状态
        hasContainer,
        showInitModel,
        initProgress,
        isInitializing,
        isCheckingStatus,
        isChangingModel,
        modelChangeProgress,
        currentModel,
        modelList,
        selectedModel,
        // 状态修改
        setShowInitModel,
        setSelectedModel,
        // 功能函数
        handleInitialize,
        checkProgress,
        checkCurrentStatus,
        getCurrentModel,
        getModelList,
        handleModelChange,
        checkModelChangeProgress,
    };
};