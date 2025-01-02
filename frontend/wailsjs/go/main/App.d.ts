// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {model} from '../models';

export function DownLoadNovel(arg1:model.SearchResult):Promise<model.CrawlResult>;

export function GetConfig():Promise<model.GetConfigResult>;

export function GetCurrentUseModel():Promise<model.GetCurrentUseModelResult>;

export function GetDownloadProgress(arg1:model.SearchResult):Promise<model.ProgressResult>;

export function GetInitOllamaProgress():Promise<model.InitOllamaProgressResult>;

export function GetSelectModelList():Promise<model.GetSelectModelListResult>;

export function GetSetOllamaModelProgress():Promise<model.GetSetOllamaModelProgressResult>;

export function GetUpdateInfo():Promise<model.GetUpdateInfoResult>;

export function GetUsageInfo():Promise<model.GetUsageInfoResult>;

export function HasInitOllama():Promise<model.HasInitOllamaResult>;

export function InitOllama():Promise<model.InitOllamaResult>;

export function SerachNovel(arg1:string):Promise<Array<model.SearchResult>>;

export function SetConfig(arg1:string):Promise<string>;

export function SetOllamaModel(arg1:string):Promise<model.SetOllamaModelResult>;

export function StartChatbot(arg1:string):Promise<model.StartChatbotResult>;
