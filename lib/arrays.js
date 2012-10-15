var util=require('util');
var obj=require('./obj');

/**
 * Computes the intersection of two arrays.
 * @param a1
 * @param a2
 * @return {Array} All elements that are in a1 and in a2 (or an empty array).
 */
exports.intersect=function(a1,a2)
{
    var a,b, c,result=[];

    if (!util.isArray(a1) || !util.isArray(a2))
        return [];

    //choose to iterate over the shorter array
    if (a2.length<a1.length)
    {
        a=a2;
        b=a1;
    }
    else
    {
        a=a1;
        b=a2;
    }

    for (var i in a)
        if (obj.equal(a[i],b[i]))
            result.push(a[i]);

    return result;
}

/**
 * Computes the differences of two arrays.
 * @param a1
 * @param a2
 * @return {Array} All elements that are not in a1 (or an empty array if they are the same).
 */
exports.diff=function(a1,a2)
{
    var a,b, result=[];

    if (!util.isArray(a1))
        return a2;

    if (!util.isArray(a2))
        return a1;

    //choose to iterate over the longer array (so that any remaining elements
    if (a2.length>a1.length)
    {
        a=a2;
        b=a1;
    }
    else
    {
        a=a1;
        b=a2;
    }

    for (var i in a)
        if (!obj.equal(a[i],b[i]))
            result.push(a[i]);

    return result;
}

exports.equal=function(a,b)
{
    return obj.equal(a,b);
}
