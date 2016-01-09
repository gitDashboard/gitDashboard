var gitDashboard = angular.module("gitDashboard",['ngRoute','angular-jwt','LocalStorageModule','authService','reposService','ui.gravatar','userService','groupService','dtrw.bcrypt']);

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
	}).
	when('/users',{
		templateUrl:"public/fragment/users.html",
		controller:"UsersController",
		controllerAs:"users"
	}).
	when('/groups',{
		templateUrl:"public/fragment/groups.html",
		controller:"GroupsController",
		controllerAs:"groups"
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

gitDashboard.controller('SelUserController',['$scope','User','$location',function($scope,User,$location){
	$scope.search=function(){
		User.search($scope.username).then(function(data){
			if (data.success){
				$scope.foundedUsers=data.users;
			}else{
				alert(data.error.message);
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.selectUser=function(user){
		$scope.currPerm.userName=user.username;
		$scope.currPerm.userId=user.id;
		$scope.currPerm.groupName =null;
		$scope.currPerm.groupId =null;
		$('#searchUserPopup').modal('hide');
	}
}]);

gitDashboard.controller('SelGroupController',['$scope','Group','$location',function($scope,Group,$location){
	$scope.search=function(){
		Group.search($scope.name).then(function(data){
			if (data.success){
				$scope.foundedGroups=data.groups;
			}else{
				alert(data.error.message);
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.selectGroup=function(group){
		$scope.currPerm.userName=null;
		$scope.currPerm.userId=null;
		$scope.currPerm.groupName =group.name;
		$scope.currPerm.groupId =group.id;
		$('#searchGroupPopup').modal('hide');
	}
}]);


