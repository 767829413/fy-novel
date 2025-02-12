import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Tabs } from 'antd';

const { TabPane } = Tabs;



const FunnyToy: React.FC = () => {
    const { t } = useTranslation();
    const [activeTab, setActiveTab] = useState('image');

    const handleTabChange = (activeKey: string) => {
        setActiveTab(activeKey);
        switch (activeKey) {
            case 'image':
                break;
        }
    };


    return (
        <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
            <Tabs activeKey={activeTab} onChange={handleTabChange}>
                <TabPane tab="Pictures" key="image">
                    <div style={{ marginBottom: '20px' }}>
                    </div>
                </TabPane>
            </Tabs>
        </div>
    );
};

export default React.memo(FunnyToy);