export namespace model {
	
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

}

