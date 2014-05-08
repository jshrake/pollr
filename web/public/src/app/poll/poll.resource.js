/*jslint browser: true*/
/*global angular*/

'use strict';

require('angular');

module.exports = angular.module('pollr.PollResource', [])
    .factory('Poll', ['$resource',
        function ($resource) {
            return $resource("api/polls/:id", {
                id: '@id'
            }, {});
        }
    ]);