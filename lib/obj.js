var util=require('util');

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
