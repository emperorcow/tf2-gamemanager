var admin = angular.module('admin', [
	'ngResource'
])

admin.factory('RconService', ['$http', function ($http) {
	var RconService = {};
	RconService.run = function(cmd) {
		return $http.post("/api/rcon", cmd);
	}
	return RconService;
}])

admin.controller('AdminController', ['$interval', '$scope', 'RconService', 'TeamService', function ScoreboardController($interval, $scope, RconService, TeamService) {
	$scope.rconoutput = "Output from command will appear here..."
	$scope.runCommand = function() {
		RconService.run($scope.rconcommandbox)
			.success(function(data) {
				$scope.rconoutput = data
			})
			.error(function(error) {
				$scope.rconoutput = "ERROR("+error.status+"): Unable to run command."
			})
	}

	$scope.red = TeamService.get({ name: "Red" })
	$scope.blue = TeamService.get({ name: "Blue" })
}]);

admin.factory('TeamService', ['$resource', function ($resource) {
	return $resource('/api/teams/:name', {name: '@name'});
}])


admin.directive("rconButton", ['RconService', function rconButton(RconService) {
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
						alert("Command successful.")
					})
					.error(function(error) {
						alert("Unable to run command.")
					})
			}
			$scope.clickRed = function() {
				var tmp = $scope.cmd.replace("%TAR%", "@RED")
				$scope.runCommand(tmp)
			}
			$scope.clickBlue = function() {
				var tmp = $scope.cmd.replace("%TAR%", "@BLUE")
				$scope.runCommand(tmp)
			}
		}
	}
}])