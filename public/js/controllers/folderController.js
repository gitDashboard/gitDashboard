gitDashboard.controller('FolderController',['$scope','$location','Folder','$routeParams',function($scope,$location,Folder,$routeParams){
	var folderId = $routeParams.folderId;
	Folder.getFolder(folderId).then(function(data){
		$scope.folder=data.folder
	},function(error){
		console.log(error);
		if (error.status==401){
			$location.path("login");
		}
	});
	$scope.selectUser=function(user){
		$scope.folder.admins.push(user);
		$('#searchUserPopup').modal('hide');
	}
	$scope.removeUser=function(index){
		$scope.folder.admins.splice(index,1);
	}
	$scope.saveAdmins=function(){
		Folder.setAdmins($scope.folder.id,$scope.folder.admins).then(function(data){
			if (!data.success){
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