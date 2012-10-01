var express=require('express');
var DB=require('./db');
var sconfig=require("./config");
var log=require('logule').init(module);
var path=require('path');
var ejs=require('ejs');

ejs.open='{{';
ejs.close='}}';

function addModule(app,name,mconf)
{
    var mname=sconfig.module_prefix+'-'+name;

    var module=require(mname);

    mconf=mergeConfig(sconfig,mconf);

    //merge app config with module config
	module.config=mergeConfig(module.config,mconf);

	//set the brand name (allows for easier change later on)
	module.config.brand=sconfig.brand;

    app.log.info('Mounting '+mname+' on '+module.config.prefix)

	//mount module
    app.use(module.config.prefix,module.server(app,express));

    return module;
}

function mergeConfig(source,config)
{
	//iterate over each key in config
	for (var p in config)
	{
		//if both keys are objects, merge objects recursively
		if ((typeof(source[p])==='object') && (typeof(config[p])==='object'))
			source[p]=mergeConfig(source[p],config[p]);
		else
			//override value in source
			source[p]=config[p];
	}

	return source;
}

function routing(app,mroutes)
{
	//path to <module>/routes/
	var route_dir=path.normalize(path.join(app.dirname,'/routes'));

	//loop over all methods: GET/POST/PUT/DELETE
	for (var m in mroutes)
	{
		var routes=mroutes[m];

		//loop over all routes defined for each method
		for (var r in routes)
		{
			var route=routes[r];
			var cname='index';
			var action='index';

			//objects have the form: controller:action
			if (typeof(route)==='object')
			{
				for (var c in route)
				{
					cname=c;
					action=route[c];
					//it doesn't make sense to have more than 1 controller/route
					break;
				}
			}
			else
				//string routes are synonyms for <controller>:'index'
				cname=route;

			//load the controller
			var controller=require(path.join(route_dir,cname));

			if (typeof(controller[action])!=='function')
			{
				app.log.error("route ["+m+" '"+r+"']:","no such action:",cname+'::'+action);
				continue;
			}

			var method=m.toLowerCase();

			//link the route to the action
			app[method](r,controller[action]);

			app.log.debug('['+m+" '"+r+"'] ->",cname+'::'+action);
		}
	}
}

function globalErrorHandler(err,req,res,next)
{
    throw err;

    res.send({
        success:0,
        message:err.toString()
    },500);
}

exports.run=function(config)
{
    var app=module.app=express.createServer();
    module.config=config;

    log.info('Waking up');

    app.configure(function(){
        app.use(express.methodOverride());
        app.use(express.bodyParser());
        app.use(app.router);
        app.log=log;
        app.error(globalErrorHandler);
    });

    app.configure('dev',function(){
        app.use(express.errorHandler({
            showStack: true,
            dumpExceptions: true
        }));
    });

    app.configure('prod',function(){
    	app.log.suppress('debug');
        app.use(express.errorHandler());
    });

	//expose utility methods
    app.mergeConfig=mergeConfig;
    app.addModule=addModule;
    app.routing=routing;
    app.db=DB;

    //root module must be added last, to prevent regex paths
    //from conflicting
    var last=null;

    //load modules
    for (var m in config.modules)
    {
        var mconf=config.modules[m];				//module config

        if (mconf.prefix==='/')						//check mount point
            last=m;									//if /, leave for last
        else
            addModule(app,m,mconf);					//add module
    }

	//load last module
    if (last)
        addModule(app,last,config.modules[last]);

	module.app.log.info('HTTP Server listening on '+module.config.host+':'+module.config.port);
	module.app.listen(module.config.port,module.config.host);
}
