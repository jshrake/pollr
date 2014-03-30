/*jslint browser: true*/
/*global WebSocket */
'use strict';

module.exports = (function () {

    function DisplayPollCtrl($scope, $location, $routeParams, $http, $log, Poll) {
        // Make a web socket connection to the correct poll
        var ws = new WebSocket("ws://localhost:12345/polls/" + $routeParams.pollId),
            that = this;

        Poll.get({
            id: $routeParams.pollId
        }).$promise.then(function (value) {
            that.question = value.question;
            that.choices = value.choices;
        });

        $scope.$on('elementClick.directive', function (event, data) {
            $log.info(event);
            $log.info(data);
            $http.put('http://localhost:8081/polls/' + $routeParams.pollId, {
                choiceId: data.index
            });
        });

        ws.onmessage = function (e) {
            var choiceID = parseInt(e.data, 10);
            $scope.$apply(function () {
                that.choices[choiceID].votes += 1;
            });
        };
    }

    DisplayPollCtrl.prototype.xFunction = function () {
        return function (d) {
            return d.text;
        };
    };

    DisplayPollCtrl.prototype.yFunction = function () {
        return function (d) {
            return d.votes + 1;
        };
    };

    return DisplayPollCtrl;

}());