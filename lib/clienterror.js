/** lib/clienterror.js
 *
 * @author Victor Petrov <victor.petrov@gmail.com>
 * @copyright (c) 2012, The Neuroinformatics Research Group at Harvard University.
 * @copyright (c) 2012, The President and Fellows of Harvard College.
 * @license New BSD License (see LICENSE file for details).
 */

var HTTP_BAD_REQUEST    =   400,
    HTTP_UNAUTHORIZED   =   401,
    HTTP_SERVER_ERROR   =   500;

function ClientError(message, code) {
    "use strict";

    this.name       = "ClientError";
    this.message    = message || "";
    this.code       = code || HTTP_BAD_REQUEST; //this code will attempt to read
}

ClientError.prototype = new Error();
ClientError.prototype.constructor = ClientError;
ClientError.prototype.HTTP_BAD_REQUEST  = HTTP_BAD_REQUEST;
ClientError.prototype.HTTP_UNAUTHORIZED = HTTP_UNAUTHORIZED;
ClientError.prototype.HTTP_SERVER_ERROR = HTTP_SERVER_ERROR;

module.exports = ClientError;
