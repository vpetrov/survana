/** config.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var path=require('path');

exports.brand="Survana";
exports.module_prefix="survana";


exports.encryption={
    'bits':1024,
    'key':"private/local.private"
}

exports.to_requirejs=function()
{
    var result={};

    for (l in this.lib)
        result[l]=path.join(this.lib[l],l);

    return result;
}
