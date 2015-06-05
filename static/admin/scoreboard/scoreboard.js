var scoreboard = angular.module('admin', [
	'ngResource',
	'angular-growl'
])

scoreboard.config(['growlProvider', function(growlProvider) {
	growlProvider.globalTimeToLive(5000);
}]);

scoreboard.controller('AdminController', ['$interval', '$scope', 'TeamService', 'ChallengeService', 'growl', function ScoreboardController($interval, $scope, TeamService, ChallengeService, growl) {
	ChallengeService.query()
		.success(function(data) {
			$scope.challenges = data
		})

	$scope.red = TeamService.get({ name: "Red" })
	$scope.blue = TeamService.get({ name: "Blue" })

	timer = $interval(function() {
		TeamService.get({ name: "Red" },
			function success(data, headers) {
				$scope.red.info.challenges = data.info.challenges
				$scope.red.info.score = data.info.score
				$scope.red.info.credits = data.info.credits
			})
		TeamService.get({ name: "Blue" },
			function success(data, headers) {
				$scope.blue.info.challenges = data.info.challenges
				$scope.blue.info.score = data.info.score
				$scope.blue.info.credits = data.info.credits
			})
	}, 1000);
}]);

scoreboard.factory('TeamService', ['$resource', function ($resource) {
	return $resource('/api/teams/:name', {name: '@name'});
}])

scoreboard.factory('ChallengeService', ['$http', function ($http) {
	var ChallengeService = {};
	ChallengeService.query = function() {
		return $http.get("/api/challenges");
	}
	ChallengeService.set = function(team, id, status) {
		var tmp = {
			"id": id,
			"team": team,
			"status": status
		}
		return $http.post("/api/challenges", tmp)
	}
	return ChallengeService;
}])