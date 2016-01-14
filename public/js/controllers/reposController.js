gitDashboard.controller('ReposController',['$scope','$location','Repo','$routeParams',function($scope,$location,Repo,$routeParams){
	if (!$scope.isLogged()){
		$location.path("login");	
	}
	if ($routeParams.path!=undefined){
		$scope.setCurrDir($routeParams.path)
	}
	$scope.repositories =[];
	
	$scope.showRepo=function(path,repo){
		if (repo!=null && repo.isRepo){
			$location.path("repo/"+repo.id);
		}else{
			$location.path("").search({path:path});
			
		}
	}

	$scope.createFolder=function(){
		if($scope.newFolderName!=null && $scope.newFolderName!=""){
			Repo.createFolder($scope.currDir+"/"+$scope.newFolderName).then(function(data){
				if (data.success){
					$('#createFolderPopup').modal('hide');
					$scope.newFolderName=null;
					$scope.list();
				}else{
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			})
		}
	};
	$scope.createRepo=function(){
		if($scope.newRepoName!=null && $scope.newRepoName!=""){
			Repo.createRepo($scope.currDir+"/"+$scope.newRepoName,$scope.newRepoDescription).then(function(data){
				if (data.success){
					$('#createRepoPopup').modal('hide');
					$scope.newRepoName=null;
					$scope.newRepoDescription=null;
					$scope.list();
				}else{
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			})
		}
	};
	$scope.initRepo=function(){
		Repo.initRepo($scope.currDir).then(function(data){
			if (data.success){
				$scope.upDir();
				alert("Repo initialized");
			}else{
				alert(data.error.message);
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		})
	};

	$scope.list=function(){
		Repo.list($scope.currDir).then(function(data){
			console.log(data);
			$scope.repositories = data.repositories;
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	};
	$scope.list();
}]);