/** db.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var mongodb=require('mongodb');
var util=require('./util');

function connect(fnSuccess,fnError)
{
	if (module.client)
		fnSuccess(module.client);

	module.db.open(function(error,client){
		if (error && fnError)
			return fnError(error);

		module.client=client;

		//setup shortcuts
		client.uniqueId=uniqueId;

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
					if (error)
					{
						if (fnError)
							return fnError(error,client);

						throw error;
					}

					//call the success callback
					return fnSuccess(col,client);
				});
			},fnError);
}

function find(options)
{

}

function uniqueId(dbcollection,field,callback)
{
		var id=util.randomId();
		var query={};
		query[field]=id;
		var fields={};
		fields[field]=1;

		//run a query
        dbcollection.findOne(query,fields,function(err,item){
        	//loop until item is not found
        	if (item)
        		uniqueId(dbcollection,field,callback);
        	else
        		callback(err,id);
        });
}

module.exports=function(config){
	module.db=new mongodb.Db(config.name,
							 new mongodb.Server(config.host,config.port,config.server_options),
							 config.db_options);

	return {
		db:module.db,
		connect:connect,
		collection:collection,
		find:find,
		uniqueId:uniqueId
	}
}

