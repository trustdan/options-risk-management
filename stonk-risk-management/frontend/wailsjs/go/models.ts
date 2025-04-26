export namespace models {
	
	export class RiskAssessment {
	    id: string;
	    date: time.Time;
	    emotionalScore: number;
	    fomoScore: number;
	    biasScore: number;
	    overallScore: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new RiskAssessment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = this.convertValues(source["date"], time.Time);
	        this.emotionalScore = source["emotionalScore"];
	        this.fomoScore = source["fomoScore"];
	        this.biasScore = source["biasScore"];
	        this.overallScore = source["overallScore"];
	        this.notes = source["notes"];
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
	export class StockRating {
	    id: string;
	    date: time.Time;
	    symbol: string;
	    sector: string;
	    stockSentiment: number;
	    priceTarget: number;
	    confidence: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new StockRating(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = this.convertValues(source["date"], time.Time);
	        this.symbol = source["symbol"];
	        this.sector = source["sector"];
	        this.stockSentiment = source["stockSentiment"];
	        this.priceTarget = source["priceTarget"];
	        this.confidence = source["confidence"];
	        this.notes = source["notes"];
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

export namespace time {
	
	export class Time {
	
	
	    static createFrom(source: any = {}) {
	        return new Time(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

