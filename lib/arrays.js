var util = require('util');
var obj = require('./obj');

/**
 * Computes the intersection of two arrays.
 * @param a1
 * @param a2
 * @return {Array} All elements that are in a1 and in a2 (or an empty array).
 */
exports.intersect = function (a1, a2) {
    var a, b, result = [];

    if (!util.isArray(a1) || !util.isArray(a2)) {
        return [];
    }

    //choose to iterate over the shorter array
    if (a2.length < a1.length) {
        a = a2;
        b = a1;
    } else {
        a = a1;
        b = a2;
    }

    for (var i in a) {
        if (a.hasOwnProperty(i)) {
            if (obj.equal(a[i], b[i])) {
                result.push(a[i]);
            }
        }
    }

    return result;
};

/**
 * Computes the differences of two arrays.
 * @param a1
 * @param a2
 * @return {Array} All elements that are not in a1 (or an empty array if they are the same).
 */
exports.diff = function (a1, a2) {
    var a, b, result = [];

    if (!util.isArray(a1)) {
        return a2;
    }

    if (!util.isArray(a2)) {
        return a1;
    }

    //choose to iterate over the longer array (so that any remaining elements
    if (a2.length > a1.length) {
        a = a2;
        b = a1;
    } else {
        a = a1;
        b = a2;
    }

    for (var i in a) {
        if (a.hasOwnProperty(i)) {
            if (!obj.equal(a[i], b[i]))
                result.push(a[i]);
        }
    }

    return result;
};

/**
 *
 * @param a
 * @param b
 * @return {*}
 */
exports.equal = function (a, b) {
    return obj.equal(a, b);
};

/**
 *
 * @param a
 * @param p
 */
exports.blacklist = function (a, p) {
    for (var i = 0; i < a.length; i += 1) {
        for (var j in p) {
            if (p.hasOwnProperty(j)) {
                if (a[i].hasOwnProperty(j)) {
                    delete a[i][j];
                }
            }
        }
    }
};

/**
 * Returns all elements in a1 that are not in a2
 * @param a1
 * @param a2
 */
exports.unique = function (a1, a2) {
    var result = [];

    for (var i = 0; i < a1.length; i += 1) {
        if (a2.indexOf(a1[i]) < 0) {
            result.push(a1[i]);
        }
    }

    return result;
}
