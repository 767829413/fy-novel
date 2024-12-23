export namespace config {
	
	export class Info {
	    // Go type: struct { SourceID int "mapstructure:\"source-id\" json:\"source-id\""; DownloadPath string "mapstructure:\"download-path\" json:\"download-path\""; Extname string "mapstructure:\"extname\" json:\"extname\""; LogLevel string "mapstructure:\"log-level\" json:\"log-level\"" }
	    base: any;
	    // Go type: struct { Threads int "mapstructure:\"threads\" json:\"threads\"" }
	    crawl: any;
	    // Go type: struct { MaxAttempts int "mapstructure:\"max-attempts\" json:\"max-attempts\"" }
	    retry: any;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base = this.convertValues(source["base"], Object);
	        this.crawl = this.convertValues(source["crawl"], Object);
	        this.retry = this.convertValues(source["retry"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace model {
	
	export class CrawlResult {
	    OutputPath: string;
	    TakeTime: number;
	
	    static createFrom(source: any = {}) {
	        return new CrawlResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OutputPath = source["OutputPath"];
	        this.TakeTime = source["TakeTime"];
	    }
	}
	export class GetConfigResult {
	    Config: config.Info;
	
	    static createFrom(source: any = {}) {
	        return new GetConfigResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Config = this.convertValues(source["Config"], config.Info);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetUpdateInfoResult {
	    UpdateInfo: string;
	    LatestUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new GetUpdateInfoResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UpdateInfo = source["UpdateInfo"];
	        this.LatestUrl = source["LatestUrl"];
	    }
	}
	export class GetUsageInfoResult {
	    Title: string;
	    VersionInfo: string;
	    Address: string;
	    CurrentBookSource: string;
	    ExportFormat: string;
	
	    static createFrom(source: any = {}) {
	        return new GetUsageInfoResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.VersionInfo = source["VersionInfo"];
	        this.Address = source["Address"];
	        this.CurrentBookSource = source["CurrentBookSource"];
	        this.ExportFormat = source["ExportFormat"];
	    }
	}
	export class ProgressResult {
	    Exists: boolean;
	    Completed: number;
	    Total: number;
	
	    static createFrom(source: any = {}) {
	        return new ProgressResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Exists = source["Exists"];
	        this.Completed = source["Completed"];
	        this.Total = source["Total"];
	    }
	}
	export class SearchResult {
	    url: string;
	    bookName: string;
	    author: string;
	    intro: string;
	    latestChapter: string;
	    latestUpdate: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.bookName = source["bookName"];
	        this.author = source["author"];
	        this.intro = source["intro"];
	        this.latestChapter = source["latestChapter"];
	        this.latestUpdate = source["latestUpdate"];
	    }
	}

}

