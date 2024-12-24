import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { GetConfig, SetConfig } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { Spin, Typography, Form, Input, Slider, Button, message, Select, Tooltip, Tabs, InputNumber } from 'antd';
import { QuestionCircleOutlined } from '@ant-design/icons';

const { Title } = Typography;
const { TabPane } = Tabs;

const ViewConfig: React.FC = () => {
    const { t } = useTranslation();
    const [config, setConfig] = useState<model.GetConfigResult | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [form] = Form.useForm();

    useEffect(() => {
        const fetchConfig = async () => {
            try {
                const result = await GetConfig();
                console.log("Fetched config:", result);
                setConfig(result);
                if (result.Config) {
                    form.setFieldsValue(result.Config);
                }
            } catch (error) {
                console.error(t('viewConfig.fetchError'), error);
                message.error(t('viewConfig.fetchErrorMessage'));
            } finally {
                setLoading(false);
            }
        };

        fetchConfig();
    }, [form, t]);

    const onFinish = async (values: model.GetConfigResult) => {
        try {
            const configJson = JSON.stringify(values, null, 2);
            const res = await SetConfig(configJson);
            if (res.length == 0) {
                message.success(t('viewConfig.updateSuccess'));
            } else {
                message.error(res);
            }
        } catch (error) {
            console.error(t('viewConfig.updateError'), error);
            message.error(t('viewConfig.updateErrorMessage'));
        }
    };

    if (loading) {
        return <Spin size="large" />;
    }

    const threadMarks = {
        '-1': '-1',
        '1': '1',
        '2': '2',
        '4': '4',
        '6': '6',
        '8': '8',
        '12': '12',
        '24': '24',
        '32': '32'
    };

    return (
        <div style={{ maxWidth: '600px', margin: '0 auto' }}>
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '20px'
            }}>
                <Title level={2} style={{ margin: 0 }}>{t('viewConfig.title')}</Title>
                <Tooltip title={t('viewConfig.updateConfigTooltip')}>
                    <Button
                        type="primary"
                        onClick={() => form.submit()}
                    >
                        {t('viewConfig.updateConfig')}
                    </Button>
                </Tooltip>
            </div>
            {config && config.Config && (
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={onFinish}
                    initialValues={config.Config}
                    style={{ maxWidth: '500px' }}
                >
                    <Tabs defaultActiveKey="1">
                        <TabPane tab={t('viewConfig.basicSettings')} key="1">
                            <Form.Item
                                name={["base", "source-id"]}
                                label={
                                    <span>
                                        {t('viewConfig.enableBookSource')}
                                        <Tooltip title={t('viewConfig.bookSourceTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                                style={{ width: '30%' }}
                            >
                                <Select
                                    options={[
                                        { value: 1, label: '1' },
                                        { value: 2, label: '2' },
                                    ]}
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>
                            <Form.Item
                                name={["base", "download-path"]}
                                label={
                                    <span>
                                        {t('viewConfig.downloadPath')}
                                        <Tooltip title={t('viewConfig.downloadPathTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                            >
                                <Input />
                            </Form.Item>
                            <Form.Item
                                name={["base", "extname"]}
                                label={
                                    <span>
                                        {t('viewConfig.fileExtension')}
                                        <Tooltip title={t('viewConfig.fileExtensionTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                                style={{ width: '30%' }}
                            >
                                <Select
                                    options={[
                                        { value: "txt", label: "txt" },
                                        { value: "epub", label: "epub" }
                                    ]}
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>
                            <Form.Item
                                name={["base", "log-level"]}
                                label={
                                    <span>
                                        {t('viewConfig.logLevel')}
                                        <Tooltip title={t('viewConfig.logLevelTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                                style={{ width: '30%' }}
                            >
                                <Select
                                    options={['panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'].map(level => ({
                                        value: level,
                                        label: level
                                    }))}
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>
                        </TabPane>
                        <TabPane tab={t('viewConfig.crawlSettings')} key="2">
                            <Form.Item
                                name={["crawl", "threads"]}
                                label={
                                    <span>
                                        {t('viewConfig.crawlThreads')}
                                        <Tooltip title={t('viewConfig.crawlThreadsTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                            >
                                <Slider
                                    marks={threadMarks}
                                    step={null}
                                    min={-1}
                                    max={32}
                                    included={false}
                                />
                            </Form.Item>
                        </TabPane>
                        <TabPane tab={t('viewConfig.retrySettings')} key="3">
                            <Form.Item
                                name={["retry", "max-attempts"]}
                                label={
                                    <span>
                                        {t('viewConfig.maxRetries')}
                                        <Tooltip title={t('viewConfig.maxRetriesTooltip')}>
                                            <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                                        </Tooltip>
                                    </span>
                                }
                            >
                                <InputNumber min={0} />
                            </Form.Item>
                        </TabPane>
                    </Tabs>
                </Form>
            )}
        </div>
    );
};

export default ViewConfig;