/** index.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var express =   require('express'),
    DB      =   require('./db'),
    sconfig =   require("./config"),
    log     =   require('logule').init(module),
    path    =   require('path'),
    ejs     =   require('ejs'),
    ursa    =   require('ursa'),
    fs      =   require('fs'),
    util    =   require('./util'),
    //constants
    HTTP_SERVER_ERROR = 500;

ejs.open    =   '{{';
ejs.close   =   '}}';

/* Globals */
global.obj          =   require('./lib/obj');
global.arrays       =   require('./lib/arrays');
global.ClientError  =   require('./lib/clienterror');
global.ROOT         =   path.dirname(process.mainModule.filename);

/**
 * Mount a new Survana module
 * @param app
 * @param name
 * @param mconf
 * @return {Object}
 */
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

/**
 * Overrides default module config with the global config
 * TODO: Consider using obj.merge() or obj.override()
 * @param source
 * @param config
 * @return {Object}
 */
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

    log.error(err.message, err.stack);

    if (req.header('Content-Type') === 'application/json') {
        res.send({
            success:0,
            message:err.message
        },err.code || HTTP_SERVER_ERROR);
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

	//already loaded?
	if (ursa.isKey(keypath))
		continue;

        if (!fs.existsSync(keypath))
            throw Error("'"+i+"': no key could be found at location '",keypath,"'");

        //read the key and store it instead of the 'key' property
        items[i].key=ursa.coercePublicKey(fs.readFileSync(keypath));
        items[i].keyID=items[i].key.toPublicSshFingerprint('hex');
    }

    return items;
}

exports.run=function(config) {
	//HTTP
	var http_server = httpServer(config);
	if (http_server) {
		module.app.log.info('HTTP Server listening on '+module.config.host+':'+module.config.port);
		http_server.listen(module.config.port,module.config.host);
	}
	
	//HTTPS
	if (config.https) {
		var https_server = httpsServer(config);
		if (https_server) {
			https_server.listen(module.config.https.port,module.config.https.host);
			module.app.log.info('HTTPS Server listening on ' + module.config.https.host + ':' + module.config.https.port);
		}
	}
}

function httpServer(config) {
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
    if (last) {
        addModule(app,last,config.modules[last]);
    }

    return app;
}

function httpsServer(config) {

	if (!fs.existsSync(config.https.key)) {
		log.error("HTTPS Key does not exist: " + config.https.key);
		return null;
	}

	if (!fs.existsSync(config.https.cert)) {
		log.error("HTTPS certificate does not exist: " + config.https.cert);
		log.info("Use this command to generate a self-signed certificate: \n\n" +
			"openssl req -x509 -new -key " + (config.https.key || "private/local.private") +
			" > " + (config.https.cert || "private/localhost.cert") + "\n");
		return null;
	}

    var app=express.createServer({
		key: fs.readFileSync(config.https.key),
		cert: fs.readFileSync(config.https.cert)
    });

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
    if (last) {
        addModule(app,last,config.modules[last]);
    }

    return app;
}
