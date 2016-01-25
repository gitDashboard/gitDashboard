var gitDashboard = angular.module("gitDashboard",['ngRoute','angular-jwt','LocalStorageModule','authService','reposService','ui.gravatar','userService','dtrw.bcrypt','eventService','folderService']);

gitDashboard.config(function($interpolateProvider) {
	$interpolateProvider.startSymbol('{[{');
	$interpolateProvider.endSymbol('}]}');
});

gitDashboard.filter('unsafe', function($sce) {
	return function(val) {
		return $sce.trustAsHtml(val);
	};
});

gitDashboard.config(function (localStorageServiceProvider) {
	localStorageServiceProvider
	.setPrefix('gitDashboard')
	.setStorageType('sessionStorage');
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
	when('/folder/admins/:folderId',{
		templateUrl:"public/fragment/folderAdmins.html",
		controller:"FolderController",
		controllerAs:"folder"
	}).
	when('/permission/folder/:folderId',{
		templateUrl:"public/fragment/permission.html",
		controller:"PermissionController",
		controllerAs:"folderCtrl"
	}).
	when('/permission/repo/:repoId',{
		templateUrl:"public/fragment/permission.html",
		controller:"PermissionController",
		controllerAs:"repoCtrl"
	});
}]);

gitDashboard.controller('MainCtrl', ['$scope','localStorageService','jwtHelper','$location','Repo','Folder','$route','$q', function ($scope,localStorageService,jwtHelper,$location,Repo,Folder,$route,$q) {
	$scope.isLogged=function(){
		var token =  localStorageService.get('jwt_token');
		if (token  != undefined){
			var expired = jwtHelper.isTokenExpired(token);
			if (expired){
				alert("session expired");
				$scope.logout();
				return false;
			}else{
				return true;
			}
		}
		return false;
	}
	$scope.getUser=function(){
		if (!$scope.isLogged()){
			return null;
		}else{
			user = jwtHelper.decodeToken( localStorageService.get('jwt_token'));
			user.IsAdmin = function(){
				return user.Admin
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

	$scope.upDir=function(){
		if ($scope.currDir.parentId>0){
			setCurrDir($scope.currDir.parentId);
		}
	}

	$scope.setCurrDir=function(folderId){
		if (folderId!=0){
			Folder.getFolder(folderId).then(function(data){
				if (data.success){
					$scope.currDir=data.folder;
				}else{
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}else{
			$scope.currDir={id:0};
		}
	}

	$scope.isFolderAdmin=function(){
		var user = $scope.getUser();
		if (user.Admin){
			return true;
		}
		if ($scope.currDir.id==0){
			return false;
		}
		var admin = false;
		if ($scope.currDir!=null && $scope.currDir.extAdmins!=null){
			for (var a =0; a<$scope.currDir.extAdmins.length && !admin; a++){
				admin = $scope.currDir.extAdmins[a].id==user.ID;
			}
		}
		return admin;
	}

	$scope.moveRepos=function(){
		var repoToMove = $scope.selectedRepos.slice()
		$scope.currentAction = "moving";
		$('#movingDimmer').dimmer('show');
		for (var i = repoToMove.length - 1; i >= 0; i--) {
			Repo.moveRepo(repoToMove[i],$scope.currDir.id).then(function(data){
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
}]);



