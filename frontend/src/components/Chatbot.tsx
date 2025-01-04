import { InfoCircleOutlined } from '@ant-design/icons';
import React, { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Collapse, Card, Input, Button, Typography, Space, message, Modal, Progress, Select, Tag, Alert } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined, DeleteOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { StartChatbot, HasInitOllama, GetInitOllamaProgress, InitOllama, GetCurrentUseModel, GetSelectModelList, SetOllamaModel, GetSetOllamaModelProgress } from '../../wailsjs/go/main/App';
import { useModelChange } from '../context/ModelChangeContext';

const { Title, Text } = Typography;
const { Option } = Select;
const { Panel } = Collapse;

interface ChatMessage {
    role: 'user' | 'assistant';
    content: string;
}

const Chatbot: React.FC = () => {
    const { t } = useTranslation();
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [inputText, setInputText] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [isLoading, setIsLoading] = useState(false);

    // 初始化状态
    const [showInitModel, setShowInitModal] = useState(false);
    const [initProgress, setInitProgress] = useState(0);
    const [isInitializing, setIsInitializing] = useState(false);
    const [isCheckingStatus, setIsCheckingStatus] = useState(true);

    // 初始化
    const handleInitialize = useCallback(async () => {
        setIsInitializing(true);
        setInitProgress(0);
        InitOllama();
    }, [t]);

    const checkProgress = useCallback(async () => {
        const result = await GetInitOllamaProgress();
        if (result.Exists) {
            const progressPercentage = Math.min(Math.round((result.Completed / result.Total) * 100), 100);
            setInitProgress(progressPercentage);
            if (progressPercentage >= 100) {
                setShowInitModal(false);
                setIsInitializing(false);
                setIsCheckingStatus(false)
            }
        }
        checkCurrentStatus()
    }, [t]);

    useEffect(() => {
        let intervalId: number | undefined;

        if (isInitializing) {
            intervalId = window.setInterval(checkProgress, 2000);
        }

        return () => {
            if (intervalId !== undefined) {
                window.clearInterval(intervalId);
            }
        };
    }, [isInitializing, checkProgress]);

    // 检查状态
    const checkCurrentStatus = useCallback(async () => {
        const result = await HasInitOllama();
        console.info(result)
        // 没有容器
        setShowInitModal(!result.Has);
        // 正在初始化
        setIsInitializing(result.IsInit);
        // 需要检查初始化状态
        setIsCheckingStatus(result.IsInit)
        // 正在更改model
        setIsChangingModel(result.IsSetModel)
    }, [t]);

    // 每次进入都要检查
    useEffect(() => {
        checkCurrentStatus();
    }, [checkCurrentStatus]);

    // 切换model
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
            checkCurrentStatus()
        };
        checkProgress();
    }, [getCurrentModel, t, setIsChangingModel, setModelChangeProgress]);

    useEffect(() => {
        getCurrentModel();
        getModelList();
    }, [getCurrentModel, getModelList]);

    const handleSendMessage = useCallback(async () => {
        if (inputText.trim() === '' || isInitializing || isChangingModel) return;

        const userMessage: ChatMessage = { role: 'user', content: inputText };
        setMessages(prevMessages => [...prevMessages, userMessage]);
        setInputText('');
        setIsLoading(true);

        try {
            const response = await StartChatbot(inputText);
            if (response.ErrorMsg) {
                throw new Error(response.ErrorMsg);
            }
            const assistantMessage: ChatMessage = { role: 'assistant', content: response.Response };
            setMessages(prevMessages => [...prevMessages, assistantMessage]);
        } catch (error) {
            console.error('Error sending message:', error);
            message.error(t('chatbot.sendError'));
        } finally {
            setIsLoading(false);
        }
    }, [inputText, isChangingModel, t]);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    const handleModelSelect = (value: string) => {
        setSelectedModel(value);
        Modal.confirm({
            title: t('chatbot.confirmModelChange'),
            content: value === currentModel
                ? t('chatbot.modelResetWarning')
                : t('chatbot.modelChangeWarning'),
            onOk() {
                handleModelChange(value);
            },
            onCancel() {
                setSelectedModel(null);
            },
        });
    };

    const handleClearChat = useCallback(() => {
        setMessages([]);
    }, []);


    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <Collapse
                ghost
                style={{ marginBottom: '20px' }}
                defaultActiveKey={['1']}
            >
                <Panel
                    header={
                        <span style={{
                            color: '#1890ff',
                            fontWeight: 'bold',
                            fontSize: '16px',
                            display: 'flex',
                            alignItems: 'center'
                        }}>
                            <InfoCircleOutlined style={{ marginRight: '8px', fontSize: '20px' }} />
                            {t('chatbot.prerequisiteTitle')}
                        </span>
                    }
                    key="1"
                    style={{
                        border: '1px solid #1890ff',
                        borderRadius: '4px',
                        backgroundColor: '#e6f7ff'
                    }}
                >
                    <div style={{ padding: '16px', backgroundColor: 'white', borderRadius: '4px' }}>
                        <Title level={4} style={{ marginBottom: '16px' }}>{t('chatbot.prerequisiteDescription')}</Title>
                        <Space direction="vertical" size="middle" style={{ width: '100%', marginBottom: '16px' }}>
                            <Card size="small" style={{ borderLeft: '4px solid #1890ff' }}>
                                <Text strong>{t('chatbot.prerequisiteDocker')}</Text>
                            </Card>
                            <Card size="small" style={{ borderLeft: '4px solid #1890ff' }}>
                                <Text strong>{t('chatbot.prerequisiteOllama')}</Text>
                            </Card>
                        </Space>
                        <Alert
                            message={t('chatbot.prerequisiteWarning')}
                            type="warning"
                            showIcon
                            style={{ marginBottom: '0' }}
                        />
                    </div>
                </Panel>
            </Collapse>
            {isCheckingStatus ? (
                <div>{t('chatbot.initialize')}...</div>
            ) : (
                <>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
                        <Title level={2}>{t('chatbot.title')}</Title>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                            <ThunderboltOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
                            <Tag color="blue" style={{ fontSize: '16px', padding: '5px 10px' }}>
                                {t('chatbot.currentModel')}
                            </Tag>
                            <Text
                                strong
                                style={{
                                    fontSize: '18px',
                                    background: 'linear-gradient(45deg, #1890ff, #52c41a)',
                                    WebkitBackgroundClip: 'text',
                                    WebkitTextFillColor: 'transparent',
                                }}
                            >
                                {currentModel}
                            </Text>
                            <Select
                                style={{ width: 200 }}
                                value={selectedModel || currentModel}
                                onChange={handleModelSelect}
                                disabled={isChangingModel}
                                showSearch
                                placeholder={t('chatbot.selectModel')}
                                filterOption={(input, option) =>
                                    option && typeof option.value === 'string'
                                        ? option.value.toLowerCase().includes(input.toLowerCase())
                                        : false
                                }
                            >
                                {modelList.map(model => (
                                    <Option key={model} value={model}>{model}</Option>
                                ))}
                            </Select>
                        </div>
                    </div>
                    {isChangingModel && (
                        <Progress percent={modelChangeProgress} status="active" />
                    )}
                    <div style={{ height: '400px', overflowY: 'auto', border: '1px solid #d9d9d9', padding: '10px', marginBottom: '20px' }}>
                        {messages.map((msg, index) => (
                            <div key={index} style={{
                                marginBottom: '10px',
                                display: 'flex',
                                flexDirection: msg.role === 'user' ? 'row-reverse' : 'row',
                                alignItems: 'flex-start',
                            }}>
                                <div style={{
                                    width: '32px',
                                    height: '32px',
                                    borderRadius: '50%',
                                    backgroundColor: msg.role === 'user' ? '#1890ff' : '#52c41a',
                                    display: 'flex',
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                    marginLeft: msg.role === 'user' ? '8px' : '0',
                                    marginRight: msg.role === 'user' ? '0' : '8px',
                                }}>
                                    {msg.role === 'user' ? <UserOutlined style={{ color: 'white' }} /> : <RobotOutlined style={{ color: 'white' }} />}
                                </div>
                                <div style={{
                                    maxWidth: '70%',
                                    padding: '8px 12px',
                                    borderRadius: '12px',
                                    backgroundColor: msg.role === 'user' ? '#e6f7ff' : '#f0f0f0',
                                }}>
                                    <Text>{msg.content}</Text>
                                </div>
                            </div>
                        ))}
                        <div ref={messagesEndRef} />
                    </div>
                    <Space.Compact style={{ width: '100%' }}>
                        <Input
                            value={inputText}
                            onChange={(e) => setInputText(e.target.value)}
                            onPressEnter={handleSendMessage}
                            placeholder={t('chatbot.inputPlaceholder')}
                            disabled={isInitializing || isLoading || isChangingModel}
                        />
                        <Button
                            type="primary"
                            icon={<SendOutlined />}
                            onClick={handleSendMessage}
                            loading={isLoading}
                            disabled={isInitializing || isLoading || isChangingModel}
                        >
                            {t('chatbot.send')}
                        </Button>
                        <Button
                            icon={<DeleteOutlined />}
                            onClick={handleClearChat}
                            disabled={messages.length === 0 || isChangingModel}
                        >
                            {t('chatbot.clear')}
                        </Button>
                    </Space.Compact>
                    <Modal
                        title={t('chatbot.initModalTitle')}
                        open={showInitModel}
                        onOk={handleInitialize}
                        onCancel={() => setShowInitModal(false)}
                        okText={t('chatbot.initModalOk')}
                        cancelText={t('chatbot.initModalCancel')}
                        confirmLoading={isInitializing}
                    >
                        <p>{t('chatbot.initModalContent')}</p>
                        {isInitializing && <Progress percent={initProgress} status="active" />}
                    </Modal>
                </>
            )}
        </div>
    );
};

export default React.memo(Chatbot);