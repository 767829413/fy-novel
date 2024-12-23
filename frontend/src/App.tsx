import React, { useState } from 'react';
import { Layout, Menu, theme } from 'antd';
import './App.css';
import { DownloadProvider } from './context/DownloadContext';
import DownloadNovel from './components/DownloadNovel';
import CheckUpdate from './components/CheckUpdate';
import ViewConfig from './components/ViewConfig';
import UsageInfo from './components/UsageInfo';

const { Header, Content, Footer } = Layout;
const options = [
    "下载小说",
    "检查更新",
    "查看配置文件",
    "使用须知",
];

const App: React.FC = () => {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    const [selectedOption, setSelectedOption] = useState<number>(0);

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
                return <p>请选择一个选项</p>;
        }
    };

    const items = options.map((option, index) => ({
        key: index.toString(),
        label: option,
    }));

    return (
        <Layout>
            <Header style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Menu
                    theme="dark"
                    mode="horizontal"
                    items={items}
                    defaultSelectedKeys={[`0`]}
                    style={{
                        flex: 'none',
                        minWidth: 0,
                        display: 'flex',
                        justifyContent: 'center',
                    }}
                    onClick={handleMenuClick}
                />
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
                我是混子 ©{new Date().getFullYear()} Created by 呵呵呵
            </Footer>
        </Layout>
    );
};

export default App;