gitDashboard.controller('GroupsController',['$scope','Group','$location',function($scope,Group,$location){
	$scope.list=function(){
		Group.list().then(function(data){
			if (data.success){
				$scope.groups=data.groups;
			}else{
				alert(data.error.message);
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	};
	$scope.addGroup= function(){
		$scope.groups.push(null);
	}
	$scope.saveGroup =  function(group){
		Group.save(group).then(function(data){
			if (!data.success){
				alert(data.error.message);
			}else{
				$scope.list();
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}	
		})
	}
	$scope.deleteGroup=  function(group){
		Group.deleteGroup(group).then(function(data){
			if (!data.success){
				alert(data.error.message);
			}else{
				$scope.list();
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}	
		})
	}
	$scope.selGroup=function(group){
		if ($scope.currGroup!=null && $scope.currGroup!=group){
			$scope.currGroup.showUser=false;
		}
		$scope.currGroup = group;
		group.showUser = !group.showUser;
	}
	$scope.removeUserFromGroup=function(userIndex){
		$scope.currGroup.users.splice(userIndex,1);
	}
	$scope.addUserToCurrentGroup=function(user){
		$scope.currGroup.users.push(user)
	}
	$scope.list();
}]);