import React, { useState, useEffect } from 'react';
import { QRCode, Space, Input, Button, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { model } from "../../wailsjs/go/models";
import { GetUpdateInfo } from "../../wailsjs/go/main/App.js"
import { CopyOutlined } from '@ant-design/icons';

const CheckUpdate: React.FC = () => {
    const { t } = useTranslation();
    const [updateInfo, setUpdateInfo] = useState<model.GetUpdateInfoResult | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchUpdateInfo = async () => {
            try {
                const result = await GetUpdateInfo();
                setUpdateInfo(result);
            } catch (error) {
                console.error(t('checkUpdate.errorFetching'), error);
                setUpdateInfo({
                    ErrorMsg: t('checkUpdate.fetchFailure'),
                    NeedUpdate: false,
                    LatestVersion: '',
                    CurrentVersion: '',
                    LatestUrl: ''
                });
            } finally {
                setLoading(false);
            }
        };

        fetchUpdateInfo();
    }, [t]);

    const copyToClipboard = () => {
        if (updateInfo?.LatestUrl) {
            navigator.clipboard.writeText(updateInfo.LatestUrl)
                .then(() => {
                    message.success(t('checkUpdate.copySuccess'));
                })
                .catch(() => {
                    message.error(t('checkUpdate.copyFailure'));
                });
        }
    };

    const renderUpdateInfo = () => {
        if (!updateInfo) return null;

        if (updateInfo.ErrorMsg) {
            return <p>{updateInfo.ErrorMsg}</p>;
        }

        if (updateInfo.NeedUpdate) {
            return (
                <>
                    <p>{t('checkUpdate.newVersionAvailable', { version: updateInfo.LatestVersion })}</p>
                    <p>{t('checkUpdate.currentVersion', { version: updateInfo.CurrentVersion })}</p>
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
                            {t('checkUpdate.copy')}
                        </Button>
                    </Space>
                </>
            );
        }

        return <p>{t('checkUpdate.upToDate')}</p>;
    };

    return (
        <div>
            <h2>{t('checkUpdate.title')}</h2>
            {loading ? (
                <p>{t('checkUpdate.checking')}</p>
            ) : (
                <Space direction="vertical" align="center" style={{ width: '100%' }}>
                    {renderUpdateInfo()}
                </Space>
            )}
        </div>
    );
};

export default CheckUpdate;