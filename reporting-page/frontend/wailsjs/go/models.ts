export namespace checks {
	
	export class Check {
	    issue_id: number;
	    result_id: number;
	    result?: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Check(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_id = source["issue_id"];
	        this.result_id = source["result_id"];
	        this.result = source["result"];
	        this.error = source["error"];
	    }
	}

}

export namespace scan {
	
	export class Severity {
	    checkid: number;
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

