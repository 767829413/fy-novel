import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Table, Space, Input, Button, message, Progress } from 'antd';
import { EventsOn } from "../../wailsjs/runtime/runtime";
import type { TableProps } from 'antd';
import { SerachNovel, DownLoadNovel, GetDownloadProgress } from "../../wailsjs/go/main/App.js"
import { model } from "../../wailsjs/go/models";
import { useDownload } from '../context/DownloadContext';

const { Search } = Input;

const DownloadNovel: React.FC = () => {
    const { t } = useTranslation();
    const [searchResults, setSearchResults] = useState<model.SearchResult[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [downloadProgress, setDownloadProgress] = useState<number>(0);
    const [isMerging, setIsMerging] = useState<boolean>(false);
    const progressIntervalRef = useRef<number | null>(null);
    const { isDownloading, setIsDownloading } = useDownload();
    const [currentPage, setCurrentPage] = useState<number>(1);
    const pageSize = 10;

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
        localStorage.removeItem('searchResults');
        try {
            const results = await SerachNovel(value);
            setSearchResults(results);
            try {
                localStorage.setItem('searchResults', JSON.stringify(results));
            } catch (storageError) {
                console.error(t('downloadNovel.storageError'), storageError);
                message.warning(t('downloadNovel.storageWarning'));
            }
        } catch (error) {
            console.error(t('downloadNovel.searchError'), error);
            message.error(t('downloadNovel.searchFailure'));
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
        message.info(t('downloadNovel.startDownload', { bookName: record.bookName }));

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
                console.error(t('downloadNovel.progressError'), error);
            }
        }, 1000);

        DownLoadNovel(record)
            .then((result: model.CrawlResult) => {
                if (progressIntervalRef.current) {
                    clearInterval(progressIntervalRef.current);
                }
                setDownloadProgress(100);
                setIsMerging(false);
                message.success(t('downloadNovel.downloadComplete', { 
                    bookName: record.bookName, 
                    outputPath: result.OutputPath, 
                    takeTime: result.TakeTime 
                }));
            })
            .catch((error) => {
                if (progressIntervalRef.current) {
                    clearInterval(progressIntervalRef.current);
                }
                console.error(t('downloadNovel.downloadError'), error);
                message.error(t('downloadNovel.downloadFailure', { bookName: record.bookName }));
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
            title: t('downloadNovel.columnIndex'),
            key: 'index',
            render: (_, __, index) => (currentPage - 1) * pageSize + index + 1,
            align: 'center',
        },
        {
            title: t('downloadNovel.columnBookName'),
            dataIndex: 'bookName',
            key: 'BookName',
            align: 'center',
        },
        {
            title: t('downloadNovel.columnAuthor'),
            dataIndex: 'author',
            key: 'Author',
            align: 'center',
        },
        {
            title: t('downloadNovel.columnLatestChapter'),
            dataIndex: 'latestChapter',
            key: 'LatestChapter',
            align: 'center',
        },
        {
            title: t('downloadNovel.columnLatestUpdate'),
            dataIndex: 'latestUpdate',
            key: 'LatestUpdate',
            align: 'center',
        },
        {
            title: t('downloadNovel.columnAction'),
            key: 'action',
            align: 'center',
            render: (_, record) => (
                <Space size="middle">
                    <Button
                        onClick={() => handleDownload(record)}
                        disabled={isDownloading}
                    >
                        {t('downloadNovel.buttonDownload')}
                    </Button>
                </Space>
            ),
        },
    ];

    return (
        <div>
            <h2>{t('downloadNovel.startSearch')}</h2>
            <Space>
                <Search
                    placeholder={t('downloadNovel.searchPlaceholder')}
                    enterButton={t('downloadNovel.buttonSearch')}
                    size="large"
                    loading={loading}
                    onSearch={handleSearch}
                    value={searchQuery}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
                />
                <Button onClick={handleClear} size="large">{t('downloadNovel.buttonClear')}</Button>
            </Space>
            <br />
            <br />
            {isDownloading && (
                <Progress
                    percent={downloadProgress}
                    status="active"
                    format={(percent: number = 0) => {
                        if (isMerging) {
                            return t('downloadNovel.merging');
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
                    total: searchResults ? searchResults.length : 0,
                    onChange: (page: number) => setCurrentPage(page),
                }}
            />
        </div>
    );
};

export default DownloadNovel;