/** index.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var express=require('express');
var DB=require('./db');
var sconfig=require("./config");
var log=require('logule').init(module);
var path=require('path');
var ejs=require('ejs');
var ursa=require('ursa');
var fs=require('fs');
var util=require('./util');

ejs.open='{{';
ejs.close='}}';

global.obj=require('./lib/obj');
global.arrays=require('./lib/arrays');
global.ROOT=path.dirname(process.mainModule.filename);

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
    mserver.keyID=app.keyID;
    mserver.randomId=util.randomId;

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

function routing(app,mconfig,custom)
{
	//path to <module>/routes/
	var mroutes     =   mconfig.routes,
        route_dir   =   path.normalize(path.join(app.dirname,'/routes')),
        middleware;

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

            //call custom function to determine middleware, if any
            if (custom) {
                middleware = custom(app, mconfig, m,r,cname,action,controller,controller[action]);
            }

            //link the route to the action
            if (middleware) {
                app[method](r,middleware,controller[action]);
            } else {
                app[method](r,controller[action]);
            }


			app.log.debug('['+m+" '"+r+"'] ->",cname+'::'+action+'   '+(middleware?'[M]':'[-]'));
		}
	}
}

function globalErrorHandler(err,req,res,next)
{
    var app=req.app;
    var log=app.log;    //use app-specific logger

    log.error(err.message,err.stack);

    if (req.header('Content-Type') === 'application/json') {
        res.send({
            success:0,
            message:err.message
        },500);
    } else {
        res.render('error',{
            req:req,
            app:app,
            err:err
        })
    }
}

function getServerKey(encryption)
{
    if (!encryption.key)
        throw Error('Configuration: No encryption key defined (encryption.key)');

    var key=null;                                                   //an instance of an 'ursa' private key
    var keyPath=path.join(ROOT,encryption.key).split('.');          //join application root dir and 'key' config value

    //remove file extension by deleting the last string after a dot
    if (keyPath.length>1)
        keyPath.pop();

    //join all strings
    keyPath=keyPath.join('.');

    var privateKeyPath=keyPath+'.private';  //append ".private" to private keys
    var publicKeyPath =keyPath+'.public';   //append ".public"  to public keys

    //if the private key has already been generated, use it - the public key can be regenerated from it
    if (fs.existsSync(privateKeyPath))
        key=ursa.coercePrivateKey(fs.readFileSync(privateKeyPath)); //load private key from disk
    else
    {
        //if not - generate a new key
        key=ursa.generatePrivateKey(encryption.bits);               //generate new key with specified # of bits

        fs.writeFileSync(privateKeyPath,key.toPrivatePem());        //save key to disk
        fs.writeFileSync(publicKeyPath,key.toPublicPem());          //save public key in PEM format

        /* Note: The public key is saved so that it could be shared with other modules for signing purposes. The idea is
                 to let the user copy 'local.public' to 'user@server:/www/private/harvard.public' (i.e. some other
                 survana component that needs to interact with this instance of survana) and then the user can point
                 a specific component (such as 'survana-study') to this public key, which would then allow this instance
                 to send signed requests to the other instance.

           TODO: Use a database to store the keys. Provide an UI for key exchanges and remote instance registration
        */
    }

    return key;
}

function readKeys(items)
{
    //load all keys
    for (var i in items)
    {
        //skip elements without a 'key' property
        if (!items[i]['key'])
            continue;

        var keypath=items[i].key;

        if (!fs.existsSync(keypath))
            throw Error("'"+i+"': no key could be found at location '"+keypath+"'");

        //read the key and store it instead of the 'key' property
        items[i].key=ursa.coercePublicKey(fs.readFileSync(keypath));
        items[i].keyID=items[i].key.toPublicSshFingerprint('hex');
    }

    return items;
}

exports.run=function(config)
{
    var app=module.app=express.createServer();
    module.config=config;

    var key=getServerKey(config.encryption);

    log.info('Waking up');

    app.configure(function(){
        app.set('views', __dirname + '/views');
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
    app.readKeys=readKeys;
    app.db=DB;
    app.publicKey=ursa.coercePublicKey(key.toPublicPem());
    app.privateKey=ursa.coercePrivateKey(key);
    app.keyID=app.publicKey.toPublicSshFingerprint('hex');

    //root module must be added last, to prevent regex paths
    //from conflicting
    var last=null;

    //load modules
    for (var m in config.modules)
    {
        var mconf=config.modules[m];				//module config

        mconf.publicURL = config.publicURL;

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
