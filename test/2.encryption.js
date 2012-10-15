var vows=require('vows');
var assert=require('assert');
var ursa=require('ursa');

var suite=vows.describe('encryption');

suite.addBatch({
   'RSA':{
       topic:ursa.generatePrivateKey(),
       'can generate keypair':function(topic){
           assert.isObject(topic);
           assert.isTrue(ursa.isKey(topic));
       },
       'can export public key':function(topic){
           assert.include(topic.toPublicPem().toString(),'PUBLIC KEY');
       },
       'can export private key':function(topic){
           assert.include(topic.toPrivatePem().toString(),'PRIVATE KEY')
       }
   }
});

suite.export(module);
