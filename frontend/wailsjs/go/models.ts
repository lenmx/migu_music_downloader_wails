export namespace model {
	
	export class BaseResponse {
	    code: number;
	    message: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new BaseResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class I18nSourceMap {
	    currentLang: string;
	    sources: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new I18nSourceMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentLang = source["currentLang"];
	        this.sources = source["sources"];
	    }
	}
	export class Setting {
	    language: string;
	    savePath: string;
	    downloadLrc: boolean;
	    downloadCover: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Setting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.savePath = source["savePath"];
	        this.downloadLrc = source["downloadLrc"];
	        this.downloadCover = source["downloadCover"];
	    }
	}

}

