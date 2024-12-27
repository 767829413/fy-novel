import React, { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Input, Button, Typography, Space } from 'antd';
import { SendOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons';
import { StartChatbot } from '../../wailsjs/go/main/App';

const { Title, Text } = Typography;

interface ChatMessage {
    text: string;
    isUser: boolean;
}

class ChatbotAPIError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "ChatbotAPIError";
    }
}

class NetworkError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "NetworkError";
    }
}

const Chatbot: React.FC = () => {
    const { t } = useTranslation();
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [inputText, setInputText] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [isLoading, setIsLoading] = useState(false);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(scrollToBottom, [messages]);

    const callChatbotAPI = async (message: string): Promise<string> => {
        try {
            const response = await StartChatbot(message);
            if (response.ErrorMsg && response.ErrorMsg.length > 0) {
                throw new ChatbotAPIError(response.ErrorMsg);
            }

            if (response.Response && response.Response.length > 0) {
                return response.Response;
            }
            throw new Error(t('chatbot.invalidResponse'));
        } catch (error) {
            console.error('Error calling chatbot API:', error);
            if (error instanceof ChatbotAPIError) {
                return t('chatbot.apiError', { error: error.message });
            } else if (error instanceof NetworkError) {
                return t('chatbot.networkError');
            } else if (error instanceof Error) {
                return t('chatbot.generalError', { error: error.message });
            } else {
                return t('chatbot.unknownError');
            }
        }
    };

    const handleSendMessage = async () => {
        if (inputText.trim() === '') return;

        const newUserMessage: ChatMessage = {
            text: inputText,
            isUser: true,
        };

        setMessages([...messages, newUserMessage]);

        setIsLoading(true);
        try {
            const botReply = await callChatbotAPI(inputText);
            const botMessage: ChatMessage = {
                text: botReply,
                isUser: false,
            };
            setMessages(prevMessages => [...prevMessages, botMessage]);
        } catch (error) {
            console.error('Error in handleSendMessage:', error);
            const errorMessage: ChatMessage = {
                text: t('chatbot.errorMessage'),
                isUser: false,
            };
            setMessages(prevMessages => [...prevMessages, errorMessage]);
        } finally {
            setIsLoading(false);
            setInputText('');
        }
    };

    return (
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
                />
                <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSendMessage}
                    loading={isLoading}
                    disabled={isLoading || inputText.trim() === ''}
                >
                    {t('chatbot.send')}
                </Button>
            </Space.Compact>
        </Space>
    );
};

export default Chatbot;