gitDashboard.controller('ReposController',['$scope','$location','Repo','Folder','$routeParams',function($scope,$location,Repo,Folder,$routeParams){
	if (!$scope.isLogged()){
		$location.path("login");	
	}

	$scope.showRepo=function(repo){
		$location.path("repo/"+repo.id);
	}
	$scope.openFolder=function(folder){
		$location.path("").search({folderId:folder.id});
	}
	$scope.upFolder=function(){
		$location.path("").search({folderId:$scope.currDir.parentId});
	}

	$scope.createFolder=function(){
		if($scope.newFolderName!=null && $scope.newFolderName!=""){
			Folder.createFolder($scope.currDir.id,$scope.newFolderName,$scope.newFolderDescription).then(function(data){
				if (data.success){
					$('#createFolderPopup').modal('hide');
					$scope.newFolderName=null;
					$scope.newFolderDescription=null;
					$scope.list($scope.currDir.id);
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
			Repo.createRepo($scope.currDir.id,$scope.newRepoName,$scope.newRepoDescription).then(function(data){
				if (data.success){
					$('#createRepoPopup').modal('hide');
					$scope.newRepoName=null;
					$scope.newRepoDescription=null;
					$scope.list($scope.currDir.id);
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

	$scope.list=function(folderId){
		folderId=parseInt(folderId);
		Folder.list(folderId).then(function(data){
			console.log(data);
			$scope.folders = data.folders;
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
		Repo.list(folderId).then(function(data){
			console.log(data);
			$scope.repositories = data.repositories;
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	};

	
	if ($routeParams.folderId!=undefined){
		$scope.setCurrDir($routeParams.folderId);
		$scope.list($routeParams.folderId);
	}else{
		$scope.setCurrDir(0);
		$scope.list(0);
	}
	
	$scope.repositories =[];
}]);