import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeRaw from 'rehype-raw'
import { EyeInvisibleOutlined, EyeOutlined } from '@ant-design/icons';
import { InfoCircleOutlined } from '@ant-design/icons';
import React, { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Collapse, Input, Button, Typography, Space, message, Modal, Progress, Select, Tag, Tabs } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined, DeleteOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { StartChatbot, DeepSeekChat } from '../../wailsjs/go/main/App';
import { useOllama } from '../hooks/useOllama';

const { Title, Text } = Typography;
const { Option } = Select;
const { Panel } = Collapse;
const { TabPane } = Tabs;

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
    const [activeTab, setActiveTab] = useState('deepseek');
    const [deepseekApiKey, setDeepseekApiKey] = useState('');
    const [showDeepseekApiKey, setShowDeepseekApiKey] = useState(false);

    const {
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
    } = useOllama();

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

    // 每次进入都要检查
    useEffect(() => {
        if (activeTab === 'ollama') {
            checkCurrentStatus();
        }
    }, [checkCurrentStatus]);

    useEffect(() => {
        let intervalId: number | undefined;

        if (isChangingModel) {
            intervalId = window.setInterval(checkModelChangeProgress, 2000);
        }

        return () => {
            if (intervalId !== undefined) {
                window.clearInterval(intervalId);
            }
        };
    }, [isChangingModel, checkModelChangeProgress]);

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
            let response;
            if (activeTab === 'ollama') {
                response = await StartChatbot(inputText);
                if (response.ErrorMsg) {
                    throw new Error(response.ErrorMsg);
                }
            } else {
                response = await DeepSeekChat(inputText, deepseekApiKey);
                console.info('DeepSeek response:', response);
                if (response.ErrorMsg) {
                    throw new Error(response.ErrorMsg);
                }
            }
            const assistantMessage: ChatMessage = { role: 'assistant', content: response.Response };
            setMessages(prevMessages => [...prevMessages, assistantMessage]);
        } catch (error) {
            console.error('Error sending message:', error);
            message.error(t('chatbot.sendError'));
        } finally {
            setIsLoading(false);
        }
    }, [inputText, isChangingModel, t, activeTab, deepseekApiKey]);

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
            onOk: async () => {
                await handleModelChange(value);
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
            <Tabs activeKey={activeTab} onChange={setActiveTab}>
                <TabPane tab="DeepSeek" key="deepseek">
                    <div style={{ marginBottom: '20px' }}>
                        <Input.Password
                            placeholder={t('chatbot.deepseekApiKey')}
                            value={deepseekApiKey}
                            onChange={(e) => setDeepseekApiKey(e.target.value)}
                            style={{ marginBottom: '10px' }}
                            iconRender={(visible) => (visible ? <EyeOutlined /> : <EyeInvisibleOutlined />)}
                            visibilityToggle={{ visible: showDeepseekApiKey, onVisibleChange: setShowDeepseekApiKey }}
                        />
                    </div>
                </TabPane>
                <TabPane tab="Ollama" key="ollama">
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
                            {/* Ollama 先决条件内容 */}
                        </Panel>
                    </Collapse>

                    {isCheckingStatus ? (
                        <div>
                            {t('chatbot.status')}...
                            {isInitializing && (
                                <Progress percent={initProgress} status="active" />
                            )}
                            {isChangingModel && (
                                <Progress percent={modelChangeProgress} status="active" />
                            )}
                        </div>
                    ) : (
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
                            <Title level={2}>{t('chatbot.title')}</Title>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                {/* Ollama 模型选择器 */}
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
                                    disabled={!hasContainer || isInitializing || isLoading || isChangingModel}
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
                    )}
                </TabPane>
            </Tabs>

            {/* 聊天消息显示区域 */}
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
                            {msg.role === 'user' ? (
                                <Text>{msg.content}</Text>
                            ) : (
                                <ReactMarkdown
                                    remarkPlugins={[remarkGfm]}
                                    rehypePlugins={[rehypeRaw]}
                                >
                                    {msg.content}
                                </ReactMarkdown>
                            )}
                        </div>
                    </div>
                ))}
                <div ref={messagesEndRef} />
            </div>

            {/* 输入区域 */}
            <Space.Compact style={{ width: '100%' }}>
                <Input
                    value={inputText}
                    onChange={(e) => setInputText(e.target.value)}
                    onPressEnter={handleSendMessage}
                    placeholder={t('chatbot.inputPlaceholder')}
                    disabled={activeTab === 'ollama' ? (!hasContainer || isInitializing || isLoading || isChangingModel) : false}
                />
                <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSendMessage}
                    loading={isLoading}
                    disabled={activeTab === 'ollama' ? (!hasContainer || isInitializing || isLoading || isChangingModel) : false}
                >
                    {t('chatbot.send')}
                </Button>
                <Button
                    icon={<DeleteOutlined />}
                    onClick={handleClearChat}
                    disabled={messages.length === 0 || (activeTab === 'ollama' && (!hasContainer || isInitializing || isLoading || isChangingModel))}
                >
                    {t('chatbot.clear')}
                </Button>
            </Space.Compact>

            {/* Ollama 初始化模态框 */}
            <Modal
                title={t('chatbot.initModalTitle')}
                open={showInitModel}
                onOk={handleInitialize}
                onCancel={() => setShowInitModel(false)}
                okText={t('chatbot.initModalOk')}
                cancelText={t('chatbot.initModalCancel')}
                confirmLoading={isInitializing}
            >
                <p>{t('chatbot.initModalContent')}</p>
                {isInitializing && <Progress percent={initProgress} status="active" />}
            </Modal>
        </div>
    );
};

export default React.memo(Chatbot);