import React, { useState, useEffect, useRef } from 'react';
import { Table, Space, Input, Button, message, Progress } from 'antd';
import { EventsOn } from "../../wailsjs/runtime/runtime";
import type { TableProps } from 'antd';
import { SerachNovel, DownLoadNovel, GetDownloadProgress } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";

const { Search } = Input;

const DownloadNovel: React.FC = () => {
    const [searchResults, setSearchResults] = useState<model.SearchResult[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [downloadProgress, setDownloadProgress] = useState<number>(0);
    const [isDownloading, setIsDownloading] = useState<boolean>(false);
    const [currentDownloadingBook, setCurrentDownloadingBook] = useState<string | null>(null);
    const progressIntervalRef = useRef<number | null>(null);

    useEffect(() => {
        // 从本地存储加载之前的搜索结果
        const savedResults = localStorage.getItem('searchResults');
        if (savedResults) {
            setSearchResults(JSON.parse(savedResults));
        }

        // 监听下载进度事件
        EventsOn("downloadProgress", (progress: number) => {
            setDownloadProgress(progress);
        });

        return () => {
            // 清理函数
            if (progressIntervalRef.current) {
                clearInterval(progressIntervalRef.current);
            }
        };
    }, []);

    const handleSearch = async (value: string) => {
        setLoading(true);
        setSearchQuery(value);
        try {
            const results = await SerachNovel(value);
            setSearchResults(results);
            localStorage.setItem('searchResults', JSON.stringify(results));
        } catch (error) {
            console.error("处理搜索时出错:", error);
            message.error("搜索失败，请稍后重试");
        } finally {
            setLoading(false);
        }
    };

    const handleClear = () => {
        setSearchResults([]);
        setSearchQuery('');
        localStorage.removeItem('searchResults');
    };

    const handleDownload = (record: model.SearchResult) => {
        setIsDownloading(true);
        setDownloadProgress(0);
        setCurrentDownloadingBook(record.bookName);
        message.info(`开始下载《${record.bookName}》`);

        progressIntervalRef.current = window.setInterval(async () => {
            try {
                // const progress = await GetDownloadProgress(record);
                // setDownloadProgress(progress);
            } catch (error) {
                console.error("获取下载进度时出错:", error);
            }
        }, 1000);

        DownLoadNovel(record)
            .then((result: model.CrawlResult) => {
                if (progressIntervalRef.current) {
                    clearInterval(progressIntervalRef.current);
                }
                setDownloadProgress(100);
                message.success(`《${record.bookName}》下载完成！输出路径：${result.OutputPath} 花费时间: ${result.TakeTime} 秒`);
            })
            .catch((error) => {
                if (progressIntervalRef.current) {
                    clearInterval(progressIntervalRef.current);
                }
                console.error("下载小说时出错:", error);
                message.error(`下载《${record.bookName}》失败，请稍后重试`);
            })
            .finally(() => {
                if (progressIntervalRef.current) {
                    clearInterval(progressIntervalRef.current);
                }
                setIsDownloading(false);
                setCurrentDownloadingBook(null);
                setDownloadProgress(0);
            });
    };

    const columns: TableProps<model.SearchResult>['columns'] = [
        {
            title: '序号',
            key: 'index',
            render: (_, __, index) => index + 1,
            align: 'center',
        },
        {
            title: '书名',
            dataIndex: 'bookName',
            key: 'BookName',
            align: 'center',
        },
        {
            title: '作者',
            dataIndex: 'author',
            key: 'Author',
            align: 'center',
        },
        {
            title: '最新章节',
            dataIndex: 'latestChapter',
            key: 'LatestChapter',
            align: 'center',
        },
        {
            title: '最后更新时间',
            dataIndex: 'latestUpdate',
            key: 'LatestUpdate',
            align: 'center',
        },
        {
            title: '操作',
            key: 'action',
            align: 'center',
            render: (_, record) => (
                <Space size="middle">
                    <Button
                        onClick={() => handleDownload(record)}
                        disabled={isDownloading}
                    >
                        下载
                    </Button>
                </Space>
            ),
        },
    ];

    return (
        <div>
            <h2>开始搜索</h2>
            <Space>
                <Search
                    placeholder="输入小说名称"
                    enterButton="搜索"
                    size="large"
                    loading={loading}
                    onSearch={handleSearch}
                    value={searchQuery}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
                />
                <Button onClick={handleClear} size="large">清空</Button>
            </Space>
            <br />
            <br />
            {isDownloading && (
                <div>
                    <p>正在下载: {currentDownloadingBook}</p>
                    <Progress percent={downloadProgress} status="active" />
                </div>
            )}
            <Table<model.SearchResult>
                columns={columns}
                dataSource={searchResults}
            />
        </div>
    );
};

export default DownloadNovel;