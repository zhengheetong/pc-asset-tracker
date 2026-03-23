export namespace main {
	
	export class PCSpecs {
	    cpu: string;
	    ramTotal: string;
	    ramModules: string;
	    disks: string;
	    serial: string;
	    tag1: string;
	    tag2: string;
	    tag3: string;
	
	    static createFrom(source: any = {}) {
	        return new PCSpecs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.ramTotal = source["ramTotal"];
	        this.ramModules = source["ramModules"];
	        this.disks = source["disks"];
	        this.serial = source["serial"];
	        this.tag1 = source["tag1"];
	        this.tag2 = source["tag2"];
	        this.tag3 = source["tag3"];
	    }
	}

}

