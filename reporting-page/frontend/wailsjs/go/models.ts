export namespace checks {
	
	export class Check {
	    id: string;
	    result?: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Check(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.result = source["result"];
	        this.error = source["error"];
	    }
	}

}

export namespace scan {
	
	export class DataBaseData {
	    id: string;
	    severity: number;
	    jsonkey: string;
	
	    static createFrom(source: any = {}) {
	        return new DataBaseData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.severity = source["severity"];
	        this.jsonkey = source["jsonkey"];
	    }
	}
	export class JsonKey {
	    id: string;
	    key: string;
	
	    static createFrom(source: any = {}) {
	        return new JsonKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.key = source["key"];
	    }
	}
	export class Severity {
	    checkid: string;
	    level: number;
	
	    static createFrom(source: any = {}) {
	        return new Severity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.checkid = source["checkid"];
	        this.level = source["level"];
	    }
	}

}

