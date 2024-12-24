import React, { useState, useEffect } from 'react';
import { GetConfig, SetConfig } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { Spin, Typography, Form, Input, Slider, Button, message, Select, Tooltip, Tabs, InputNumber } from 'antd';
import { QuestionCircleOutlined } from '@ant-design/icons';

const { Title } = Typography;
const { TabPane } = Tabs;

const ViewConfig: React.FC = () => {
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
                console.error("获取配置信息时出错:", error);
                message.error("获取配置信息失败");
            } finally {
                setLoading(false);
            }
        };

        fetchConfig();
    }, [form]);

    const onFinish = async (values: model.GetConfigResult) => {
        try {
            const configJson = JSON.stringify(values, null, 2);
            const res = await SetConfig(configJson);
            if (res.length == 0) {
                message.success('配置更新成功');
            } else {
                message.error(res);
            }
        } catch (error) {
            console.error("更新配置时出错:", error);
            message.error('更新配置失败');
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
        <div style={{ maxWidth: '600px', margin: '0 auto', position: 'relative' }}>
            <Title level={2} style={{ marginBottom: '20px' }}>配置一览</Title>
            <Tooltip title="建议使用默认配置以获得最佳性能">
                <Button
                    type="primary"
                    onClick={() => form.submit()}
                    style={{
                        position: 'absolute',
                        top: '0',
                        right: '0'
                    }}
                >
                    更新配置
                </Button>
            </Tooltip>
            {config && config.Config && (
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={onFinish}
                    initialValues={config.Config}
                    style={{ maxWidth: '500px' }}
                >
                    <Tabs defaultActiveKey="1">
                        <TabPane tab="基础设置" key="1">
                            <Form.Item
                                name={["base", "source-id"]}
                                label={
                                    <span>
                                        启用书源
                                        <Tooltip title="只能选一个, 当前可选值：1、2、3，建议使用默认书源 3">
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
                                        { value: 3, label: '3' }
                                    ]}
                                    style={{ width: '100%' }}
                                />
                            </Form.Item>
                            <Form.Item
                                name={["base", "download-path"]}
                                label={
                                    <span>
                                        下载路径
                                        <Tooltip title="绝对相对均可 (Windows 路径分隔符不要用 \ , 用 / 或 \)">
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
                                        文件扩展名
                                        <Tooltip title="支持 txt, epub, 推荐 epub">
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
                                        日志级别
                                        <Tooltip title="默认 error (panic fatal error warn info debug trace)">
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
                        <TabPane tab="爬取设置" key="2">
                            <Form.Item
                                name={["crawl", "threads"]}
                                label={
                                    <span>
                                        爬取线程数
                                        <Tooltip title="-1 表示自动设置，设置过大会引发防爬机制">
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
                        <TabPane tab="重试设置" key="3">
                            <Form.Item
                                name={["retry", "max-attempts"]}
                                label={
                                    <span>
                                        最大重试次数
                                        <Tooltip title="针对首次下载失败的章节">
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