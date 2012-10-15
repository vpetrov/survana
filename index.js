var express=require('express');
var DB=require('./db');
var sconfig=require("./config");
var log=require('logule').init(module);
var path=require('path');
var ejs=require('ejs');
var ursa=require('ursa');
var fs=require('fs');

ejs.open='{{';
ejs.close='}}';

global.obj=require('./lib/obj');
global.arrays=require('./lib/arrays');

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

    var mserver=module.server(app,express); //get access to the module's 'app' instance
    mserver.use(globalErrorHandler);        //register global error handler
    mserver.publicKey=app.publicKey;        //transfer server public key
    mserver.privateKey=app.privateKey;      //transfer server private key

	//mount module
    app.use(module.config.prefix,mserver);

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
    var app=req.app;
    var log=app.log;    //use app-specific logger

    log.error(err.message);

    res.send({
        success:0,
        message:err.message
    },500);
}

function getServerKey(encryption)
{
    var key=null;
    var keypath=path.join(module.parent.dirname,encryption.key);

    //if the key has already been generated, use it; if not - generate new key
    if (fs.exists(keypath))
        key=ursa.coerceKey(fs.readFileSync(keypath));
    else
    {
        key=ursa.generatePrivateKey(encryption.bits);               //generate new key
        fs.writeFileSync(keypath,key);                              //save key to disk
        fs.writeFileSync(keypath+'.pem',key.toPublicPem());         //save public key in PEM format
    }

    return key;
}

exports.run=function(config)
{
    var app=module.app=express.createServer();
    module.config=config;

    var key=getServerKey(config.encryption);

    log.info('Waking up');

    app.configure(function(){
        app.use(express.methodOverride());
        app.use(express.bodyParser());
        app.use(app.router);
        app.use(globalErrorHandler);

        app.log=log;
    });

	//expose utility methods
    app.mergeConfig=mergeConfig;
    app.addModule=addModule;
    app.routing=routing;
    app.db=DB;
    app.publicKey=ursa.coercePublicKey(key.toPublicPem());
    app.privateKey=ursa.coercePrivateKey(key);

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
