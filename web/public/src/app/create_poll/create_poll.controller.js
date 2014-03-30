'use strict';

module.exports = (function () {

    function CreatePollCtrl($location, Poll) {
        this.nextID = 0;
        this.poll = {
            'question': '',
            'choices': [],
        };
        this.newChoice = '';
        this.location = $location;
        this.Poll = Poll;
    }

    CreatePollCtrl.prototype.addChoice = function () {
        if (this.newChoice) {
            this.poll.choices.push(this.newChoice);
            this.nextID += 1;
            this.newChoice = '';
        }
    };

    CreatePollCtrl.prototype.deleteChoice = function (choice) {
        this.poll.choices.splice(this.poll.choices.indexOf(choice), 1);
    };

    CreatePollCtrl.prototype.createPoll = function () {
        var poll = new this.Poll(),
            that = this;
        this.addChoice();
        poll.question = this.poll.question;
        poll.choices = this.poll.choices;
        poll.$save().then(function () {
            that.location.path(poll.id);
            that.poll.question = '';
            that.poll.choices = [];
        });
    };

    return CreatePollCtrl;

}());