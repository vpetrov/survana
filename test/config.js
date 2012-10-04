var vows=require('vows');
var assert=require('assert');

var suite=vows.describe('config');

suite.addBatch({
    'config':{
        topic:require('../config.js'),
        'has brand name':function(topic){
            assert.isString(topic['brand']);
        },
        'has module prefix':function(topic){
            assert.isString(topic['module_prefix']);
        },
        'can convert to requirejs':function(topic){
            assert.isFunction(topic['to_requirejs']);
        }
    }
});

suite.export(module);
