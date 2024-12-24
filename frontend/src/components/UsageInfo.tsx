import React, { useState, useEffect } from 'react';
import { GetUsageInfo } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { Space, Spin, Typography, Descriptions, Tag } from 'antd';
import { RocketOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;


const UsageInfo: React.FC = () => {
    const [usageInfo, setUsageInfo] = useState<model.GetUsageInfoResult | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchUsageInfo = async () => {
            try {
                const result = await GetUsageInfo();
                setUsageInfo(result);
            } catch (error) {
                console.error("获取使用须知信息时出错:", error);
                setUsageInfo(null);
            } finally {
                setLoading(false);
            }
        };

        fetchUsageInfo();
    }, []);

    if (loading) {
        return <Spin size="large" />;
    }

    if (!usageInfo) {
        return <Paragraph>无法获取使用须知信息</Paragraph>;
    }

    return (
        <Space
            direction="vertical"
            size={16}
            style={{
                display: 'flex',
                padding: '24px',
            }}
        >
            <Title level={2}>{usageInfo.Title}</Title>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <RocketOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
                <Tag color="blue" style={{ fontSize: '16px', padding: '5px 10px' }}>
                    当前版本
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
                    {usageInfo.VersionInfo}
                </Text>
            </div>
            <Descriptions bordered column={1}>
                <Descriptions.Item label="免责声明">此程序乃作者研习Go语言之练习项目，倘使用中有何问题，皆与作者无关！！！</Descriptions.Item>
                <Descriptions.Item label="官方地址">{usageInfo.Address}</Descriptions.Item>
                <Descriptions.Item label="当前书源">{usageInfo.CurrentBookSource}</Descriptions.Item>
                <Descriptions.Item label="导出格式">{usageInfo.ExportFormat}</Descriptions.Item>
                <Descriptions.Item label="注意事项">有些不支持梯子访问哦！！！</Descriptions.Item>
            </Descriptions>
        </Space>
    );
};

export default UsageInfo;