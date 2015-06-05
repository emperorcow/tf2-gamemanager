var admin = angular.module('admin', [
	'ngResource',
	'angular-growl'
])

admin.config(['growlProvider', function(growlProvider) {
	growlProvider.globalTimeToLive(5000);
}]);

admin.factory('RconService', ['$http', function ($http) {
	var RconService = {};
	RconService.run = function(cmd) {
		return $http.post("/api/rcon", cmd);
	}
	return RconService;
}])

admin.controller('AdminController', ['$interval', '$scope', 'RconService', 'TeamService', 'ChallengeService', 'growl', function ScoreboardController($interval, $scope, RconService, TeamService, ChallengeService, growl) {
	$scope.rconoutput = "Output from command will appear here..."
	$scope.runCommand = function() {
		RconService.run($scope.rconcommandbox)
			.success(function(data) {
				growl.success("Command successfully completed.")
				$scope.rconoutput = data
			})
			.error(function(error) {
				$scope.rconoutput = ""
				growl.error("ERROR("+error.status+"): Unable to run command.")
			})
	}

	$scope.setChallenge = function(team, challenge, status) {
		ChallengeService.set(team.name, challenge.id, status)
			.success(function() {
				growl.success("Challenge successfully altered.")
				team.info.challenges[challenge.id] = status
				if(status) {
					team.info.score += challenge.value * 10
					team.info.credits += challenge.value
				} else {
					team.info.score -= challenge.value * 10
					team.info.credits -= challenge.value
				}
			})
			.error(function() {
				growl.error("Unable to alter challenge.")
			})
	}

	$scope.$watch("commonCommands", function(newVal, oldVal) {
		$scope.rconcommandbox = newVal;
	});

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
	}, 5000);
}]);

admin.factory('TeamService', ['$resource', function ($resource) {
	return $resource('/api/teams/:name', {name: '@name'});
}])

admin.factory('ChallengeService', ['$http', function ($http) {
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

admin.directive("rconButton", ['RconService', 'growl', function rconButton(RconService, growl) {
	return {
		restrict: 'E',
		replace: true,
		scope: {
			icon: '@',
			name: '@',
			cmd: '@'
		},
		template: '<div class="row action"> <div class="col-sm-6"> <h5><i class="fa fa-{{icon}}"></i> {{name}}</h5> </div> <div class="col-sm-6"> <div class="btn-group pull-right"> <button class="btn btn-danger" ng-click="clickRed()">Red</button><button class="btn btn-info" ng-click="clickBlue()">Blue</button> </div> </div> </div>',
		controller: function($scope) {
			$scope.runCommand = function(cmd) {
				RconService.run(cmd)
					.success(function(data) {
						growl.success("Rcon command successfully executed.")
					})
					.error(function(error) {
						growl.error("Unable to run command.")
					})
			}
			$scope.clickRed = function() {
				var tmp = $scope.cmd.replace("%TAR%", "@red")
				$scope.runCommand(tmp)
			}
			$scope.clickBlue = function() {
				var tmp = $scope.cmd.replace("%TAR%", "@blue")
				$scope.runCommand(tmp)
			}
		}
	}
}])