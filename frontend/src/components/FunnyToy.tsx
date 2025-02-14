import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Tabs, Input, Button, message, Space, Typography, Image, InputNumber, Switch, Collapse, Tooltip } from 'antd';
import { CopyOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { GenerateAsciiImage } from '../../wailsjs/go/main/App';

const { Text } = Typography;
const { Panel } = Collapse;

const IMAGE_TAB: string = 'image';

const FunnyToy: React.FC = () => {
    const { t } = useTranslation();
    const [activeTab, setActiveTab] = useState(IMAGE_TAB);
    const [imagePath, setImagePath] = useState('');
    const [processedImagePath, setProcessedImagePath] = useState('');
    const [threshold, setThreshold] = useState(140);
    const [filename, setFilename] = useState('tmp.txt');
    const [tmpImgName, setTmpImgName] = useState('');
    const [ascWidth, setAscWidth] = useState(-1);
    const [ascHeight, setAscHeight] = useState(-1);
    const [useYukkuri, setUseYukkuri] = useState(false);
    const [isProcessing, setIsProcessing] = useState(false);

    const handleTabChange = (activeKey: string) => {
        setActiveTab(activeKey);
    };

    const handleImagePathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setImagePath(e.target.value);
    };

    const handleSubmit = async () => {
        if (imagePath) {
            setIsProcessing(true);
            try {
                const resultPath = await GenerateAsciiImage({
                    ImgPath: imagePath,
                    Threshold: threshold,
                    Filename: filename,
                    TmpImgName: tmpImgName,
                    AscWidth: ascWidth,
                    AscHeight: ascHeight,
                    Ykr: useYukkuri,
                });
                setProcessedImagePath(resultPath.Response);
                message.success(t('funny.imageSuccessfully'));
            } catch (error) {
                console.error(t('funny.imageWhileError'), error);
                message.error(t('funny.imageError'));
            } finally {
                setIsProcessing(false);
            }
        } else {
            message.warning(t('funny.imagePathInput'));
        }
    };

    const handleCopy = () => {
        navigator.clipboard.writeText(processedImagePath).then(() => {
            message.success(t('funny.copySuccess'));
        }, () => {
            message.error(t('funny.copyError'));
        });
    };

    const items = [
        {
            key: 'image',
            label: t('funny.imageTab'),
            children: (
                <Space direction="vertical" style={{ width: '100%', marginBottom: '20px' }}>
                    <Space.Compact style={{ width: '100%' }}>
                        <Input
                            placeholder={t('funny.imagePathInput')}
                            value={imagePath}
                            onChange={handleImagePathChange}
                            style={{ width: 'calc(100% - 120px)' }}
                            disabled={isProcessing}
                        />
                        <Button
                            type="primary"
                            onClick={handleSubmit}
                            disabled={!imagePath || isProcessing}
                            style={{ width: '120px' }}
                            loading={isProcessing}
                        >
                            {isProcessing ? t('funny.processing') : t('funny.imageProcess')}
                        </Button>
                    </Space.Compact>
                    <Collapse>
                        <Panel header={t('funny.advancedSettings')} key="1">
                            <Space direction="vertical" style={{ width: '100%' }}>
                                <Tooltip title={t('funny.threshold')} placement="topLeft">
                                    <InputNumber
                                        placeholder={t('funny.threshold')}
                                        value={threshold}
                                        onChange={(value) => setThreshold(value !== null ? value : 140)}
                                        style={{ width: '100%' }}
                                        suffix={<InfoCircleOutlined style={{ color: 'rgba(0,0,0,.45)' }} />}
                                    />
                                </Tooltip>
                                <Tooltip title={t('funny.filename')} placement="topLeft">
                                    <Input
                                        placeholder={t('funny.filename')}
                                        value={filename}
                                        onChange={(e) => setFilename(e.target.value)}
                                        style={{ width: '100%' }}
                                        suffix={<InfoCircleOutlined style={{ color: 'rgba(0,0,0,.45)' }} />}
                                    />
                                </Tooltip>
                                <Tooltip title={t('funny.ascWidth')} placement="topLeft">
                                    <InputNumber
                                        placeholder={t('funny.ascWidth')}
                                        value={ascWidth}
                                        onChange={(value) => setAscWidth(value !== null ? value : -1)}
                                        style={{ width: '100%' }}
                                        suffix={<InfoCircleOutlined style={{ color: 'rgba(0,0,0,.45)' }} />}
                                    />
                                </Tooltip>
                                <Tooltip title={t('funny.ascHeight')} placement="topLeft">
                                    <InputNumber
                                        placeholder={t('funny.ascHeight')}
                                        value={ascHeight}
                                        onChange={(value) => setAscHeight(value !== null ? value : -1)}
                                        style={{ width: '100%' }}
                                        suffix={<InfoCircleOutlined style={{ color: 'rgba(0,0,0,.45)' }} />}
                                    />
                                </Tooltip>
                                <Space>
                                    <Text>{t('funny.useYukkuri')}:</Text>
                                    <Switch
                                        checked={useYukkuri}
                                        onChange={(checked) => setUseYukkuri(checked)}
                                    />
                                </Space>
                            </Space>
                        </Panel>
                    </Collapse>
                    {processedImagePath && (
                        <>
                            <Space>
                                <Text>{t('funny.processedImagePath')}:</Text>
                                <Text code>{processedImagePath}</Text>
                                <Button
                                    icon={<CopyOutlined />}
                                    onClick={handleCopy}
                                    type="link"
                                >
                                    {t('funny.copy')}
                                </Button>
                            </Space>
                            <Image
                                src={`file://${processedImagePath}`}
                                alt={t('funny.processedImage')}
                                style={{ maxWidth: '100%', maxHeight: '400px' }}
                            />
                        </>
                    )}
                </Space>
            ),
        },
    ];

    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <Tabs activeKey={activeTab} onChange={handleTabChange} items={items} />
        </div>
    );
};

export default React.memo(FunnyToy);