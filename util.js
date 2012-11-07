var alphabet='ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789'; //removed confusing letters O,o,l,I
var alphabet_max=alphabet.length-1;

/** Generates a unique alphanumeric sequence of characters
 * @param $length Number of characters to generate
 * @return String
 */
exports.randomId=function(len)
{
    if ((typeof(len)==='undefined') || !len)
        len=6;

    var result="";

    //build a string, char by char, with random characters from the 'alphabet' string
    for (var i=0;i<len;++i)
        result+=alphabet[parseInt(Math.random()*10000) % alphabet_max];

    return result;
}
