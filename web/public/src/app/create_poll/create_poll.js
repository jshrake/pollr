/*jslint browser: true*/
/*global angular*/

'use strict';

require('angular');

module.exports = angular.module('pollr.CreatePollModule', [
    require('./../poll/poll.resource.js').name
]).controller('CreatePollCtrl', ['$location', 'Poll', require('./create_poll.controller.js')]);