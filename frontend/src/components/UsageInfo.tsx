import React from 'react';

// Simulate fetching the "使用须知" content
const fetchUsageInfo = () => {
    return {
        version: "v0.0.2",
        commit: "228f345",
        buildDate: "2024-12-19",
        officialUrl: "https://github.com/767829413/easy-novel",
        sourceUrl: "http://www.mcmssc.la/",
        sourceId: 3,
        exportFormat: "epub",
    };
};

const UsageInfo: React.FC = () => {
    const usageInfo = fetchUsageInfo();

    return (
        <div>
            <h2>使用须知</h2>
            <p>fy-novel {usageInfo.version} (commit {usageInfo.commit}, built at {usageInfo.buildDate})</p>
            <p>官方地址：{usageInfo.officialUrl}</p>
            <p>当前书源：{usageInfo.sourceUrl} (ID: {usageInfo.sourceId})</p>
            <p>导出格式：{usageInfo.exportFormat}</p>
        </div>
    );
};

export default UsageInfo;