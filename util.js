/** Generates a unique alphanumeric sequence of characters
 * @param $length Number of characters to generate
 * @return String
 */
exports.randomId=function(len)
{
	 if (typeof(len)==='undefined')
	 	len=6;

     var result="";
     var alphabet='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
     var max=alphabet.length-1;

     for (var i=0;i<len;++i)
        result+=alphabet[parseInt(Math.random()*10000) % max];
     
     return result;
 }