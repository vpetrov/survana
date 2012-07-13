var prefix="idata";

exports.run=function(config)
{
	//load all modules
	for (var m in config.modules)
	{
		var mconf=config.modules[m];
		
		//do not load modules deployed on other servers	
		if (!mconf.url)
		{
			var name=prefix+'-'+m;
			var module=require(name);
			
			module.load(mconf);
			exports[m]=module;
		}
	}
};
