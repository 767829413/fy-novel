import './i18n';
import React, { useState } from 'react';
import { Layout, Menu, theme } from 'antd';
import { useTranslation } from 'react-i18next';
import './App.css';
import { DownloadProvider } from './context/DownloadContext';
import DownloadNovel from './components/DownloadNovel';
import CheckUpdate from './components/CheckUpdate';
import ViewConfig from './components/ViewConfig';
import UsageInfo from './components/UsageInfo';
import LanguageSwitcher from './components/LanguageSwitcher';

const { Header, Content, Footer } = Layout;

const App: React.FC = () => {
    const { t } = useTranslation();
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    const [selectedOption, setSelectedOption] = useState<number>(0);

    const options = [
        { key: '0', label: t('app.menu.downloadNovel') },
        { key: '1', label: t('app.menu.checkUpdate') },
        { key: '2', label: t('app.menu.viewConfig') },
        { key: '3', label: t('app.menu.usageInfo') },
    ];

    const handleMenuClick = (e: any) => {
        const index = parseInt(e.key, 10);
        setSelectedOption(index);
    };

    const renderContent = () => {
        switch (selectedOption) {
            case 0:
                return (
                    <DownloadProvider>
                        <DownloadNovel />
                    </DownloadProvider>
                );
            case 1:
                return <CheckUpdate />;
            case 2:
                return <ViewConfig />;
            case 3:
                return <UsageInfo />;
            default:
                return <p>{t('app.menu.selectOption')}</p>;
        }
    };

    return (
        <Layout>
            <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Menu
                    theme="dark"
                    mode="horizontal"
                    items={options}
                    defaultSelectedKeys={['0']}
                    style={{
                        flex: 1,
                        minWidth: 0,
                    }}
                    onClick={handleMenuClick}
                />
                <LanguageSwitcher />
            </Header>
            <Content style={{ padding: '0 48px' }}>
                <div
                    style={{
                        background: colorBgContainer,
                        minHeight: 650,
                        padding: 24,
                        borderRadius: borderRadiusLG,
                    }}
                >
                    {renderContent()}
                </div>
            </Content>
            <Footer style={{ textAlign: 'center' }}>
                {t('app.footer', { year: new Date().getFullYear() })}
            </Footer>
        </Layout>
    );
};

export default App;