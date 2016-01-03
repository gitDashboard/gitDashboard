var gitDashboard = angular.module("gitDashboard",['ngRoute','angular-jwt','LocalStorageModule','authService','reposService','ui.gravatar']);

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
	when('/login',{
		templateUrl:"public/fragment/login.html",
		controller:"LoginController",
		controllerAs:"login"
	}).
	when('/repo/:repoId',{
		templateUrl:"public/fragment/repo.html",
		controller:"RepoController",
		controllerAs:"repo"
	});
}]);

gitDashboard.controller('MainCtrl', ['$scope','localStorageService','jwtHelper','$location', function ($scope,localStorageService,jwtHelper,$location) {
	$scope.isLogged=function(){
		return localStorageService.get('jwt_token') != undefined;
	}
	$scope.getUser=function(){
		if (!$scope.isLogged()){
			return null;
		}else{
			user = jwtHelper.decodeToken( localStorageService.get('jwt_token'));
			user.IsAdmin = function(){
				return user.Groups.indexOf("admin")!=-1;
			}
			return user;
		}
	}
	$scope.logout = function () {
		localStorageService.remove('jwt_token');
		$location.path("login");
	};
}]);

gitDashboard.controller('LoginController', ['$scope','Auth','localStorageService','$location', function ($scope,Auth,localStorageService,$location) {
	$scope.types = ["internal","ldap"];
	$scope.login = function(){
		Auth.login($scope.username,$scope.password,$scope.type).then(function(data){
			if (!data.success){
				alert("Login failed");
			}else{
				localStorageService.set('jwt_token',data.jwt_token);
				$location.path("");
			}
		},function(error){
			alert(error);
		});
	};
}]);

gitDashboard.controller('ReposController',['$scope','$location','Repo','$routeParams',function($scope,$location,Repo,$routeParams){
	if ($routeParams.path!=undefined){
		$scope.currDir=$routeParams.path;
	}else{
		$scope.currDir="";
	}
	$scope.hasParent=function(){
		return $scope.currDir!=""
	}
	$scope.upDir=function(){
		slashPos = $scope.currDir.lastIndexOf("/");
		if (slashPos>-1){
			$location.path("").search({path:$scope.currDir.substring(0,slashPos)});
		}else{
			$location.path("").search({path:""});
		}
	}
	$scope.repositories =[];
	$scope.showRepo=function(path,repo){
		if (repo!=null && repo.isRepo){
			$location.path("repo/"+repo.id);
		}else{
			$location.path("").search({path:path});
			
		}
	}
	Repo.list($scope.currDir).then(function(data){
		console.log(data);
		$scope.repositories = data.repositories;
	},function(error){
		console.log(error);
		if (error.status==401){
			$location.path("login");
		}
	});
}]);

gitDashboard.controller('RepoController',['$scope','$routeParams','Repo','$location',function($scope,$routeParams,Repo,$location){
	var repoId = parseInt($routeParams.repoId);
	$scope.page = 1;
	$scope.count = 10;
	Repo.info(repoId).then(function(data){
		if (data.success){
			console.log(data);
			$scope.repo = data.info;
		}else{
			console.log(data.error);
			alert(data.error.message);
		}
	},function(error){
		console.log(error);
		alert(error);
	});

	$scope.returnToFolder = function(){
		$location.path("").search({path:$scope.repo.folderPath});
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
		Repo.commits(repoId,null,($scope.page-1)*parseInt($scope.count),parseInt($scope.count),$scope.ascending).then(function(data){
			console.log(data);
			if (data.success){
				$scope.commits = data.commits;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.getFiles=function(parent){
		Repo.files(repoId,null,parent).then(function(data){
			console.log(data);
			if (data.success){
				$scope.files = data.files;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.openFile=function(file){
		if (file.isDir){
			$scope.getFiles(file.id);
		}
		console.log(file);
	}
	$scope.getCommits();
	$scope.getFiles(null);
}]);	

