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
	    Config: string;
	    Error: string;
	
	    static createFrom(source: any = {}) {
	        return new GetConfigResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Config = source["Config"];
	        this.Error = source["Error"];
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
