import React, { useState, useEffect, useRef } from 'react';
import { Table, Space, Input, Button, message, Progress } from 'antd';
import { EventsOn } from "../../wailsjs/runtime/runtime";
import type { TableProps } from 'antd';
import { SerachNovel, DownLoadNovel, GetDownloadProgress } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { useDownload } from '../context/DownloadContext';

const { Search } = Input;

const DownloadNovel: React.FC = () => {
    const [searchResults, setSearchResults] = useState<model.SearchResult[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [downloadProgress, setDownloadProgress] = useState<number>(0);
    const [isMerging, setIsMerging] = useState<boolean>(false);
    const progressIntervalRef = useRef<number | null>(null);
    const { isDownloading, setIsDownloading } = useDownload();
    const [currentPage, setCurrentPage] = useState<number>(1);
    const pageSize = 10; // 每页显示的条目数

    useEffect(() => {
        const savedResults = localStorage.getItem('searchResults');
        if (savedResults) {
            setSearchResults(JSON.parse(savedResults));
        }

        EventsOn("downloadProgress", (progress: number) => {
            setDownloadProgress(progress);
        });

        return () => {
            if (progressIntervalRef.current) {
                clearInterval(progressIntervalRef.current);
            }
        };
    }, []);

    const handleSearch = async (value: string) => {
        setLoading(true);
        setSearchQuery(value);
        localStorage.removeItem('searchResults'); // 清空之前的搜索结果
        try {
            const results = await SerachNovel(value);
            setSearchResults(results);
            try {
                localStorage.setItem('searchResults', JSON.stringify(results));
            } catch (storageError) {
                console.error("存储搜索结果时出错:", storageError);
                message.warning("无法存储搜索结果，可能是存储空间已满");
            }
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
        setIsMerging(false);
        message.info(`开始下载《${record.bookName}》`);

        progressIntervalRef.current = window.setInterval(async () => {
            try {
                const progress = await GetDownloadProgress(record);
                if (progress.Exists) {
                    const percentage = progress.Total > 0
                        ? Math.min(99, Math.max(0, Math.floor((progress.Completed / progress.Total) * 99)))
                        : 0;
                    setDownloadProgress(percentage);

                    if (percentage === 99 && !isMerging) {
                        setIsMerging(true);
                    }
                }
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
                setIsMerging(false);
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
                setDownloadProgress(0);
                setIsMerging(false);
            });
    };

    const columns: TableProps<model.SearchResult>['columns'] = [
        {
            title: '序号',
            key: 'index',
            render: (_, __, index) => (currentPage - 1) * pageSize + index + 1,
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
                <Progress
                    percent={downloadProgress}
                    status="active"
                    format={(percent: number = 0) => {
                        if (isMerging) {
                            return "正在合成小说文件...";
                        }
                        return `${percent}%`;
                    }}
                />
            )}
            <Table<model.SearchResult>
                columns={columns}
                dataSource={searchResults}
                pagination={{
                    current: currentPage,
                    pageSize: pageSize,
                    total: searchResults ? searchResults.length : 0, // Ensure searchResults is not null
                    onChange: (page: number) => setCurrentPage(page),
                }}
            />
        </div>
    );
};

export default DownloadNovel;