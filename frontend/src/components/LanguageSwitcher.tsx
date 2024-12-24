import React from 'react';
import { useTranslation } from 'react-i18next';
import { Select } from 'antd';

const { Option } = Select;

const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return (
    <Select defaultValue={i18n.language} style={{ width: 120 }} onChange={changeLanguage}>
      <Option value="en">English</Option>
      <Option value="zh">中文</Option>
    </Select>
  );
};

export default LanguageSwitcher;