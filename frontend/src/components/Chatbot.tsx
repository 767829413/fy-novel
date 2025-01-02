import { InfoCircleOutlined } from '@ant-design/icons';
import React, { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Collapse, Card, Input, Button, Typography, Space, message, Modal, Progress, Select, Tag, Alert } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined, DeleteOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { StartChatbot } from '../../wailsjs/go/main/App';
import { useOllamaInitialization } from '../hooks/useOllamaInitialization';
import { useModelManagement } from '../hooks/useModelManagement';

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

    const {
        isInitialized,
        showInitModal,
        initProgress,
        isInitializing,
        isCheckingStatus,
        hasContainer,
        handleInitialize,
        setShowInitModal,
    } = useOllamaInitialization();

    const {
        currentModel,
        modelList,
        selectedModel,
        isChangingModel,
        modelChangeProgress,
        setSelectedModel,
        handleModelChange,
    } = useModelManagement();

    const handleSendMessage = useCallback(async () => {
        if (inputText.trim() === '' || !isInitialized || isChangingModel) return;

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
    }, [inputText, isInitialized, isChangingModel, t]);

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
                <div>正在检查状态...</div>
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
                            disabled={!isInitialized || isLoading || isChangingModel}
                        />
                        <Button
                            type="primary"
                            icon={<SendOutlined />}
                            onClick={handleSendMessage}
                            loading={isLoading}
                            disabled={!isInitialized || isLoading || isChangingModel}
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
                        open={showInitModal}
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