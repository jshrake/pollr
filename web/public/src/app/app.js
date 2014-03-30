/*jslint browser: true*/
/*global angular*/

'use strict';

require('angular');
require('angular-route');
require('angular-resource');

angular.module('pollrApp', [
    'ngRoute',
    'ngResource',
    require('./create_poll/create_poll.js').name,
    require('./display_poll/display_poll.js').name,
])
    .config(function ($routeProvider, $locationProvider, $httpProvider) {
        $routeProvider
            .when('/', {
                templateUrl: '/app/create_poll/create_poll.tpl.html'
            })
            .when('/polls/:pollId', {
                templateUrl: '/app/display_poll/display_poll.tpl.html'
            })
            .otherwise({
                redirectTo: '/'
            });
        $locationProvider.html5Mode(true);
        $httpProvider.defaults.withCredentials = true;
    });