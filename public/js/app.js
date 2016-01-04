var gitDashboard = angular.module("gitDashboard",['ngRoute','ui','angular-jwt','LocalStorageModule','authService','reposService','ui.gravatar']);

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
		Repo.commits(repoId,$scope.currRef,($scope.page-1)*parseInt($scope.count),parseInt($scope.count),$scope.ascending).then(function(data){
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
		if (parent==null){
			$scope.currPath=[];
			parentId=null;
		}else{
			parentId = parent.id;
		}
		$scope.currPath.push(parent);
		console.log($scope.currPath);
		Repo.files(repoId,$scope.currRef,parentId).then(function(data){
			console.log(data);
			if (data.success){
				$scope.parentTreeId = data.parentTreeId;
				$scope.files = data.files;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.inFolder=function(){
		return $scope.currPath.length>1;
	}
	$scope.upFilesDir=function(){
		if ($scope.currPath.length>0){
			$scope.currPath.pop();//current
			parent = $scope.currPath.pop();
			$scope.getFiles(parent);
			$scope.fileContent=null;
		}
	}
	$scope.getPath=function(){
		var strPath = "";
		for (var i=0;i<$scope.currPath.length;i++){
			if ($scope.currPath[i]!=null){
				strPath+="/"+$scope.currPath[i].name;
			}
		}
		if ($scope.file!=null){
			strPath+="/"+$scope.file.name;
		}
		return strPath;
	}

	 $scope.codemirrorLoaded = function(_editor){
	 	console.log(_editor);
	 }

	$scope.cmOption ={
		lineNumbers: true,
		matchBrackets: true,
		readOnly:true
	};

	$scope.openFile=function(file){
		if (file.isDir){
			$scope.getFiles(file);
			$scope.file=null;
			$scope.showFile=false;
		}else{
			$scope.showFile=true;
			Repo.fileContent(repoId,file.id).then(function(data){
				console.log(data);
				if (data.success){
					$scope.file={
						content:atob(data.content),
						name:file.name
					};
					$scope.fileContent=$scope.file.content;
					console.log($scope.file);
					var codeMirrorInstance = $('.CodeMirror')[0].CodeMirror;
					if (file.name.indexOf("xml")>-1){
						codeMirrorInstance.setOption("mode","text/xml")
					}
					if (file.name.indexOf("java")>-1){
						codeMirrorInstance.setOption("mode","text/x-java")
					}
					
				}else{
					alert(data.error.message)
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
		console.log(file);
	}
	$scope.showFile=false;
	$scope.currRef="refs/heads/master";	
	$scope.currPath=[];
	$scope.getCommits();
	$scope.getFiles(null);
}]);	

