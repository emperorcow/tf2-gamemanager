var team = angular.module('team', [
	'ngResource',
  'angular-growl'
])

team.factory('ActionService', ['$http', function ($http) {
	var ActionService = {};
	ActionService.query = function() {
		return $http.get("/api/actions");
	}
	ActionService.run = function(name, self) {
    var target
    if(self) {
      target = "Red"
    } else {
      target = "Red"
    }
		return $http.get("/api/actions/"+name+"/"+target)
	}
	return ActionService;
}])

team.factory('TeamService', ['$resource', function ($resource) {
	return $resource('/api/teams/:name', {name: '@name'});
}])

team.factory('ChallengeService', ['$http', function ($http) {
  var ChallengeService = {};
  ChallengeService.query = function() {
    return $http.get("/api/challenges");
  }
  return ChallengeService;
}])

team.controller('TeamController', ['$interval', '$scope', 'ActionService', 'TeamService', 'ChallengeService', 'growl', function ScoreboardController($interval, $scope, ActionService, TeamService, ChallengeService, growl) {
	ActionService.query()
		.success(function(data) {
			$scope.actions = data;
		})
		.error(function(error) {
			growl.error("An error occured gathering action data: " + error.status)	
		});

  $scope.purchaseClick = function(action, target) {
    ActionService.run(action.name, target)
      .success(function(data) {
        growl.success(action.name + " purchased successfully.")
        $scope.team.info.credits -= action.cost
      }) 
      .error(function(error) {
        growl.error("ERROR("+error.status+"): " + action.name + " could not be purchased.")
      })
  }

  ChallengeService.query()
    .success(function(data) {
      $scope.challenges = data
    })

	$scope.team = TeamService.get({ name: "Red" })

  timer = $interval(function() {
    TeamService.get({ name: "Red" },
      function success(data, headers) {
        $scope.team.info.challenges = data.info.challenges
        $scope.team.info.score = data.info.score
        $scope.team.info.credits = data.info.credits
      })
  }, 2000);
}]);