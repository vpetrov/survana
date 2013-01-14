/** lib/obj.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var util = require('util');
var UDOT='\uFF0E';
var UDOLLARSIGN='\uFF04';

exports.keys=function(o)
{
    var result=[];

    for (var i in o)
        result.push(i);

    return result;
}

/** Removes the specified properties of an object
 * @param o {Object}        The object
 * @param p {Array|Object}  Properties to remove
 * @return {Object}         The object (not a copy) without the specified properties
 * @warning The object passed as the first parameter will be modified in-place.
 */
exports.blacklist = function (o, p) {
    for (var i in p) {
        if (p.hasOwnProperty(i))
            delete o[i];
    }

    return o;
}

exports.override=function(obj1,obj2)
{
    if ((typeof obj1 === 'undefined') ||
        (typeof obj2 === 'undefined'))
        return obj1||obj2;

    //override properties of obj1
    for (var i in obj2)
        obj1[i]=obj2[i];

    return obj1;
}

exports.extract=function(obj,property,default_value,callback)
{
    if (typeof(obj[property])==='undefined')
        return default_value;

    var result=obj[property];
    delete obj[property];

    if (typeof(callback)==='function')
        callback.call(this,result,obj);

    return result;
}

exports.equal=function(obj1,obj2)
{
    //if one of them is not an object, perform strict equality
    if ((typeof(obj1)!=='object') || (typeof(obj2)!=='object'))
        return (obj1===obj2);

    var a=util.isArray(obj1);
    var b=util.isArray(obj2);

    //one of them isn't an array?
    if (a ^ b)
        return false;
    //if both are arrays or objects
    else
    {
        //arrays with different lengths?
        if (a && b && (a.length!==b.length))
            return false;

        //recursively compare these objects/arrays
        for (var i in a)
            if (!this.equal(a[i],b[i]))
                return false;
    }

    //all elements are equal, so these must be equal
    return true;
}

/**
 * Escapes all object keys containing '.' or '$' with their Unicode equivalents.
 * @param obj1 JSON object to escape
 */
exports.mongoEscape = function (obj1) {
    var n, s;

    if (typeof(obj1) !== "object")
        return;

    //loop over all elements in an array
    if (util.isArray(obj1))
    {
        n = obj1.length;

        for (var i = 0; i < n; i += 1) {
            this.mongoEscape(obj1[i]);
        }
    }
    else
    {
        for (var i in obj1) {
            if (obj1.hasOwnProperty(i)) {

                //first check the value of this key
                this.mongoEscape(obj1[i]);

                //then change the key, if necessary
                if ((typeof(i) === "string") && (i.indexOf('.') > -1)) {
                    s = i.replace(/\./g,UDOT);
                    s = s.replace(/\$/g,UDOLLARSIGN);
                    obj1[s] = obj1[i];  //copy old value to the new key
                    delete obj1[i];     //remove the old key
                }
            }
        }
    }
}
