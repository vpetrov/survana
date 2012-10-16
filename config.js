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
