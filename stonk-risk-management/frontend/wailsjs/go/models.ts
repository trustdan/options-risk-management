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
	    enthusiasm: number;
	    chartPattern: string;
	
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
	        this.enthusiasm = source["enthusiasm"];
	        this.chartPattern = source["chartPattern"];
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
	export class Trade {
	    id: string;
	    symbol: string;
	    sector: string;
	    strategy: string;
	    type: string;
	    week: number;
	    entryDate: time.Time;
	    expirationDate: time.Time;
	    entryPrice: number;
	    notes: string;
	    legNumber: number;
	    isMultiLeg: boolean;
	    shortLegExp: string;
	
	    static createFrom(source: any = {}) {
	        return new Trade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.symbol = source["symbol"];
	        this.sector = source["sector"];
	        this.strategy = source["strategy"];
	        this.type = source["type"];
	        this.week = source["week"];
	        this.entryDate = this.convertValues(source["entryDate"], time.Time);
	        this.expirationDate = this.convertValues(source["expirationDate"], time.Time);
	        this.entryPrice = source["entryPrice"];
	        this.notes = source["notes"];
	        this.legNumber = source["legNumber"];
	        this.isMultiLeg = source["isMultiLeg"];
	        this.shortLegExp = source["shortLegExp"];
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

