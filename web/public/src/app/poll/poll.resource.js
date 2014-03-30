/*jslint browser: true*/
/*global angular*/

'use strict';

require('angular');

module.exports = angular.module('pollr.PollResource', [])
    .factory('Poll', ['$resource',
        function ($resource) {
            return $resource("http://localhost:8081/polls/:id", {
                id: '@id'
            }, {});
        }
    ]);