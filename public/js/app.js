var gitDashboard = angular.module("gitDashboard",['ngRoute','angular-jwt','LocalStorageModule','authService','reposService','ui.gravatar','userService','groupService','dtrw.bcrypt','eventService']);

gitDashboard.config(function($interpolateProvider) {
	$interpolateProvider.startSymbol('{[{');
	$interpolateProvider.endSymbol('}]}');
});

gitDashboard.filter('unsafe', function($sce) {
	return function(val) {
		return $sce.trustAsHtml(val);
	};
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

gitDashboard.controller('MainCtrl', ['$scope','localStorageService','jwtHelper','$location','Repo','$route', function ($scope,localStorageService,jwtHelper,$location,Repo,$route) {
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

	$scope.selRepoIndexOf=function(repo){
		for (var i = $scope.selectedRepos.length - 1; i >= 0; i--) {
			if ($scope.selectedRepos[i].id==repo.id){
				return i;
			}
		};
		return -1;
	}

	$scope.selRepoContains=function(repo){
		return $scope.selRepoIndexOf(repo)!=-1;
	}

	$scope.selRepo=function(repo){
		var selRepoPos  =$scope.selRepoIndexOf(repo)
		if (selRepoPos==-1){
			$scope.selectedRepos.push(repo);
		}else{
			$scope.selectedRepos.splice(selRepoPos,1);
		}
	}
	
	$scope.unSelRepo=function(repo){
		var selRepoPos =$scope.selRepoIndexOf(repo)
		if (selRepoPos!=-1){
			$scope.selectedRepos.splice(selRepoPos,1);
		}
		if ($scope.currentAction == "moving" && $scope.selectedRepos.length==0 ){
			$('#movingDimmer').dimmer('hide');
			$route.reload();
		}
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
	$scope.setCurrDir=function(newDir){
		$scope.currDir=newDir;	
	}
	$scope.moveRepos=function(){
		var repoToMove = $scope.selectedRepos.slice()
		$scope.currentAction = "moving";
		$('#movingDimmer').dimmer('show');
		for (var i = repoToMove.length - 1; i >= 0; i--) {
			Repo.moveRepo(repoToMove[i],$scope.currDir).then(function(data){
				$scope.unSelRepo(data.repo);
				if (!data.success){
					console.log(data);
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
	}
	$scope.copyToClipboard=function(content){
		
	}
	$scope.selectedRepos=Array();
	$scope.currDir="";
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


