import React, { useState, useEffect } from 'react';
import { GetUpdateInfo } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { QRCode,Space, Input, Button, message } from 'antd';
import { CopyOutlined } from '@ant-design/icons';

const CheckUpdate: React.FC = () => {
    const [updateInfo, setUpdateInfo] = useState<model.GetUpdateInfoResult | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchUpdateInfo = async () => {
            try {
                const result = await GetUpdateInfo();
                setUpdateInfo(result);
            } catch (error) {
                console.error("获取更新信息时出错:", error);
                setUpdateInfo({
                    UpdateInfo: "获取更新信息失败，请稍后重试。",
                    LatestUrl: ""
                });
            } finally {
                setLoading(false);
            }
        };

        fetchUpdateInfo();
    }, []);

    const copyToClipboard = () => {
        if (updateInfo?.LatestUrl) {
            navigator.clipboard.writeText(updateInfo.LatestUrl)
                .then(() => {
                    message.success('链接已复制到剪贴板');
                })
                .catch(() => {
                    message.error('复制失败，请手动复制');
                });
        }
    };

    return (
        <div>
            <h2>检查更新</h2>
            {loading ? (
                <p>正在检查更新...</p>
            ) : updateInfo ? (
                <Space direction="vertical" align="center" style={{ width: '100%' }}>
                    <p>{updateInfo.UpdateInfo}</p>
                    {updateInfo.LatestUrl && (
                        <>
                            <QRCode value={updateInfo.LatestUrl || '-'} />
                            <Space>
                                <Input
                                    value={updateInfo.LatestUrl}
                                    readOnly
                                    style={{ width: '300px' }}
                                />
                                <Button
                                    icon={<CopyOutlined />}
                                    onClick={copyToClipboard}
                                >
                                    复制
                                </Button>
                            </Space>
                        </>
                    )}
                </Space>
            ) : (
                <p>无法获取更新信息</p>
            )}
        </div>
    );
};

export default CheckUpdate;