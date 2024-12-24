import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { GetUsageInfo } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { Space, Spin, Typography, Descriptions, Tag } from 'antd';
import { RocketOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

const UsageInfo: React.FC = () => {
    const { t } = useTranslation();
    const [usageInfo, setUsageInfo] = useState<model.GetUsageInfoResult | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchUsageInfo = async () => {
            try {
                const result = await GetUsageInfo();
                setUsageInfo(result);
            } catch (error) {
                console.error(t('usageInfo.fetchError'), error);
                setUsageInfo(null);
            } finally {
                setLoading(false);
            }
        };

        fetchUsageInfo();
    }, [t]);

    if (loading) {
        return <Spin size="large" />;
    }

    if (!usageInfo) {
        return <Paragraph>{t('usageInfo.unavailable')}</Paragraph>;
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
            <Title level={2}>{t('usageInfo.title')}</Title>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <RocketOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
                <Tag color="blue" style={{ fontSize: '16px', padding: '5px 10px' }}>
                    {t('usageInfo.currentVersion')}
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
                <Descriptions.Item label={t('usageInfo.disclaimer')}>{t('usageInfo.disclaimerContent')}</Descriptions.Item>
                <Descriptions.Item label={t('usageInfo.officialAddress')}>{usageInfo.Address}</Descriptions.Item>
                <Descriptions.Item label={t('usageInfo.currentBookSource')}>{usageInfo.CurrentBookSource}</Descriptions.Item>
                <Descriptions.Item label={t('usageInfo.exportFormat')}>{usageInfo.ExportFormat}</Descriptions.Item>
                <Descriptions.Item label={t('usageInfo.notes')}>{t('usageInfo.notesContent')}</Descriptions.Item>
            </Descriptions>
        </Space>
    );
};

export default UsageInfo;