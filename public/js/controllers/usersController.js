gitDashboard.controller('UsersController',['$scope','User','$location','bcrypt',function($scope,User,$location,bcrypt){
	$scope.list=function(){
		User.list().then(function(data){
			if (data.success){
				$scope.users=data.users;
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
	$scope.editUser=function(user){
		$scope.currUser=user;
		$('#editUserPopup').modal('show');
	}
	$scope.searchLdap=function(){
		if ($scope.currUser!=null && $scope.currUser.type=="ldap"){
			User.ldapSearch($scope.currUser.username).then(function(data){
				$scope.currUser.name = data.cn
				$scope.currUser.email = data.mail
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
	}
	
	$scope.saveUser =  function(user){
		if(user.password!=null && user.type=="internal"){
			var salt = bcrypt.genSaltSync(10);
			user.password = bcrypt.hashSync(user.password, salt);
		}
		User.save(user).then(function(data){
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
	$scope.deleteUser =  function(user){
		User.deleteUser(user).then(function(data){
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
	$scope.list();
}]);