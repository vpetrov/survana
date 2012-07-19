var mongodb=require('mongodb');

function connect(fnSuccess,fnError)
{	
	if (module.client)
		fnSuccess(module.client);
	
	module.db.open(function(error,client){
		if (error && fnError)
			return fnError(error);
			
		module.client=client;
			
		if (fnSuccess)
			return fnSuccess(client);
	});
}

function collection(name,fnSuccess,fnError)
{
	//try to connect
	connect(function(client)
			{
				//attempt to open collection by name
				return client.collection(name,function(error,col)
				{
					//check for errors
					if (error && fnError)
						return fnError(error,client);
					
					//call the success callback
					return fnSuccess(col,client);
				});
			},fnError);
}

function find(options)
{
	
}

module.exports=function(config){
	module.db=new mongodb.Db(config.name,
							 new mongodb.Server(config.host,config.port,config.server_options),
							 config.db_options);
	
	return {
		db:module.db,
		connect:connect,
		collection:collection,
		find:find
	}
}

