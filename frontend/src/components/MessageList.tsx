import React from 'react';
import { List, Typography } from 'antd';
import { UserOutlined, RobotOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
}

interface MessageListProps {
  messages: ChatMessage[];
}

const MessageList: React.FC<MessageListProps> = ({ messages }) => {
  return (
    <List
      itemLayout="horizontal"
      dataSource={messages}
      renderItem={(message) => (
        <List.Item>
          <List.Item.Meta
            avatar={message.role === 'user' ? <UserOutlined /> : <RobotOutlined />}
            title={message.role === 'user' ? 'You' : 'Assistant'}
            description={<Text>{message.content}</Text>}
          />
        </List.Item>
      )}
    />
  );
};

export default MessageList;