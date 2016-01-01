var gitDashboard = angular.module("gitDashboard",['ngRoute','angular-jwt','LocalStorageModule','authService','reposService']);

gitDashboard.config(function($interpolateProvider) {
	$interpolateProvider.startSymbol('{[{');
	$interpolateProvider.endSymbol('}]}');
});


gitDashboard.config(function ($httpProvider, jwtInterceptorProvider) {
	function localStorageService(localStorageService) {
		return localStorageService.get('jwt_token');
	}

	jwtInterceptorProvider.tokenGetter = ['localStorageService', localStorageService];
	$httpProvider.interceptors.push('jwtInterceptor');
});

gitDashboard.config(['$routeProvider', '$locationProvider',function ($routeProvider,$locationProvider) {
	$routeProvider.
	when('/',{
		templateUrl:"public/fragment/repos.html",
		controller:"ReposController",
		controllerAs:"repos"
	}).
	when('/repo/:repoId',{
		templateUrl:"public/fragment/repo.html",
		controller:"RepoController",
		controllerAs:"repo"
	});
}]);

gitDashboard.controller('MainCtrl', ['$scope','localStorageService', function ($scope,localStorageService) {
	$scope.isLogged=function(){
		return localStorageService.get('jwt_token') != undefined;
	}
}]);

gitDashboard.controller('LoginController', ['$scope','Auth','localStorageService', function ($scope,Auth,localStorageService) {
	$scope.types = ["internal","ldap"];
	$scope.login = function(){
		Auth.login($scope.username,$scope.password,$scope.type).then(function(data){
			if (!data.success){
				alert("Login failed");
			}else{
				localStorageService.set('jwt_token',data.jwt_token);
			}
		},function(error){
			alert(error);
		});
	};

	$scope.logout = function () {
		localStorageService.remove('jwt_token');
	};

}]);

gitDashboard.controller('ReposController',['$scope','$location','Repo',function($scope,$location,Repo){
	$scope.currDir="";
	$scope.repositories =[];
	$scope.showRepo=function(path,repo){
		if (repo!=null && repo.isRepo){
			$location.path("repo/"+repo.id);
		}else{
			$scope.currDir=path;
			Repo.list($scope.currDir).then(function(data){
				console.log(data);
				$scope.repositories = data.repositories;
			},function(error){
				alert(error);
			});
		}
	}
	$scope.showRepo($scope.currDir,null);
}]);

gitDashboard.controller('RepoController',['$scope','$routeParams','Repo',function($scope,$routeParams,Repo){
	var repoId = parseInt($routeParams.repoId);
	$scope.page = 1;
	$scope.count = 2;
	$scope.ascending = true;
	$scope.changeSort = function(){
		$scope.ascending=!$scope.ascending;
		$scope.getCommits();
	}
	$scope.decPage=function(){
		if ($scope.page>1){
			$scope.page--;
		}
		$scope.getCommits();
	}
	$scope.incPage=function(){
		if ($scope.commits.length>0){
			$scope.page++;
			$scope.getCommits();
		}
	}
	$scope.getCommits=function(){
		Repo.commits(repoId,null,$scope.page-1,$scope.count,$scope.ascending).then(function(data){
			console.log(data);
			if (data.success){
				$scope.commits = data.commits;
			}
		},function(error){
			alert(error);
		});
	}
	$scope.getCommits();
}]);	

