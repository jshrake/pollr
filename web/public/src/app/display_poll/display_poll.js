/*jslint browser: true*/
/*global angular*/

'use strict';

require('angular');
require('d3');
require('nvd3');
require('angularjs-nvd3-directives');

module.exports = angular.module('pollr.DisplayPollModule', [
	require('./../poll/poll.resource.js').name,
	'nvd3ChartDirectives'
]).controller('DisplayPollCtrl', ['$scope', '$location', '$routeParams', '$http', '$log', 'Poll', require('./display_poll.controller.js')]);