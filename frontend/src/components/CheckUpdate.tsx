import React from 'react';

// Simulate fetching the "检查更新" content
const fetchCheckUpdate = () => {
    return {
        result: "v0.0.2 已是最新版本！(https://github.com/767829413/easy-novel/releases/tag/v0.0.2)",
    };
};

const CheckUpdate: React.FC = () => {
    const checkUpdateResult = fetchCheckUpdate();

    return (
        <div>
            <h2>检查更新</h2>
            <p>{checkUpdateResult.result}</p>
        </div>
    );
};

export default CheckUpdate;