/* shutdown.js - connects to the MongoDB server and issues a shutdown command
    TODO: support custom host/user/password
 */

var mongo = new Mongo();
db = mongo.getDB("admin");
db.shutdownServer();