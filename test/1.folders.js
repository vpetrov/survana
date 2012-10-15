var vows=require('vows');
var assert=require('assert');
var fs=require('fs');
var path=require('path');

var suite=vows.describe('folders');

suite.addBatch({
    'private/':{
        topic:path.join(__dirname,'../private'),
        'is writable':function(topic){
            var file=path.join(topic,'.test');
            fs.writeFileSync(file,'test');
            fs.unlinkSync(file);
        }
    }
});

suite.export(module);
