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

