import React, { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Input, Button, Typography, Space, message, Modal, Progress } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons';
import { StartChatbot, GetInitOllamaProgress, InitOllama, HasInitOllama } from '../../wailsjs/go/main/App';

const { Title, Text } = Typography;

interface ChatMessage {
    text: string;
    isUser: boolean;
}

interface ChatbotResponse {
    Response: string;
    ErrorMsg: string;
}

const Chatbot: React.FC = () => {
    const { t } = useTranslation();
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [inputText, setInputText] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [isInitialized, setIsInitialized] = useState(false);
    const [showInitModal, setShowInitModal] = useState(false);
    const [initProgress, setInitProgress] = useState(0);
    const [isInitializing, setIsInitializing] = useState(false);

    useEffect(() => {
        checkInitialization();
    }, []);

    const checkInitialization = async () => {
        try {
            const result = await HasInitOllama();
            if (result.Has) {
                setIsInitialized(true);
            } else {
                setShowInitModal(true);
            }
        } catch (error) {
            console.error('Error checking Ollama initialization:', error);
            message.error(t('chatbot.checkInitError'));
        }
    };

    const handleInitialize = async () => {
        setIsInitializing(true);
        try {
            await InitOllama();
            const checkProgress = setInterval(async () => {
                const result = await GetInitOllamaProgress();
                if (result.Exists) {
                    const progressPercentage = (result.Completed / result.Total) * 100;
                    setInitProgress(progressPercentage);
                    if (progressPercentage >= 100) {
                        setIsInitialized(true);
                        setShowInitModal(false);
                        clearInterval(checkProgress);
                    }
                }
            }, 1000);
        } catch (error) {
            console.error('Error initializing Ollama:', error);
            message.error(t('chatbot.initError'));
        } finally {
            setIsInitializing(false);
        }
    };

    const callChatbotAPI = useCallback(async (message: string): Promise<string> => {
        try {
            const response: ChatbotResponse = await StartChatbot(message);
            if (response.ErrorMsg) {
                throw new Error(response.ErrorMsg);
            }
            return response.Response || t('chatbot.invalidResponse');
        } catch (error) {
            console.error('Error calling chatbot API:', error);
            return t('chatbot.apiError', { error: (error as Error).message });
        }
    }, [t]);

    const handleSendMessage = useCallback(async () => {
        if (inputText.trim() === '') return;

        const newUserMessage: ChatMessage = { text: inputText, isUser: true };
        setMessages(prevMessages => [...prevMessages, newUserMessage]);

        setIsLoading(true);
        try {
            const botReply = await callChatbotAPI(inputText);
            const botMessage: ChatMessage = { text: botReply, isUser: false };
            setMessages(prevMessages => [...prevMessages, botMessage]);
        } catch (error) {
            console.error('Error in handleSendMessage:', error);
            const errorMessage: ChatMessage = { text: t('chatbot.errorMessage'), isUser: false };
            setMessages(prevMessages => [...prevMessages, errorMessage]);
        } finally {
            setIsLoading(false);
            setInputText('');
        }
    }, [inputText, callChatbotAPI, t]);

    return (
        <>
            <Modal
                title={t('chatbot.initializationRequired')}
                open={showInitModal}
                onOk={handleInitialize}
                onCancel={() => {}}
                footer={[
                    <Button key="submit" type="primary" loading={isInitializing} onClick={handleInitialize}>
                        {t('chatbot.initialize')}
                    </Button>
                ]}
                closable={false}
                maskClosable={false}
            >
                <p>{t('chatbot.initializationDescription')}</p>
                {isInitializing && <Progress percent={initProgress} status="active" />}
            </Modal>

            <Space direction="vertical" size="large" style={{ width: '100%' }}>
                <Title level={2}>{t('chatbot.title')}</Title>
                <div style={{
                    height: '400px',
                    border: '1px solid #d9d9d9',
                    borderRadius: '4px',
                    padding: '10px',
                    overflowY: 'auto',
                    display: 'flex',
                    flexDirection: 'column'
                }}>
                    {messages.map((msg, index) => (
                        <div key={index} style={{
                            alignSelf: msg.isUser ? 'flex-end' : 'flex-start',
                            backgroundColor: msg.isUser ? '#e6f7ff' : '#f0f0f0',
                            borderRadius: '10px',
                            padding: '8px 12px',
                            marginBottom: '8px',
                            maxWidth: '70%'
                        }}>
                            <Space>
                                {msg.isUser ? <UserOutlined /> : <RobotOutlined />}
                                <Text>{msg.text}</Text>
                            </Space>
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
                        disabled={!isInitialized}
                    />
                    <Button
                        type="primary"
                        icon={<SendOutlined />}
                        onClick={handleSendMessage}
                        loading={isLoading}
                        disabled={!isInitialized || isLoading || inputText.trim() === ''}
                    >
                        {t('chatbot.send')}
                    </Button>
                </Space.Compact>
            </Space>
        </>
    );
};

export default React.memo(Chatbot);