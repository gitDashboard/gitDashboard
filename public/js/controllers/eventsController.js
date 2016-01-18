gitDashboard.controller('EventsController',['$scope','$location','Event','$routeParams',function($scope,$location,Event,$routeParams){
	$scope.page=1;
	$scope.query={
		first:0,
		count:10,
		levels:[]
	}
	if ($scope.repo!=null){
		$scope.query.repoId=$scope.repo.id;
	}
	$scope.searchEvents = function(){
		console.log($scope.query);
		$scope.query.count  = parseInt($scope.query.count);
		$scope.query.first = ($scope.page-1)*$scope.query.count ;
		Event.search($scope.query).then(function(data){
			if (!data.success){
				alert(data.error.message)
			}else{
				$scope.events = data.events;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}	
		});
	}
	$scope.incPage=function(){
		$scope.page++;
		$scope.searchEvents();
	}
	$scope.decPage=function(){
		if ($scope.page>0){
			$scope.page--;
			$scope.searchEvents();
		}
	}
	$scope.searchEvents();
}]);