	function createTooltip(stage, x,y,message){
		var length = message.length*15;
		var rect = new createjs.Shape();
		rect.graphics.setStrokeStyle(1).beginStroke("Black").beginFill("LightBlue").drawRect(x+10, y+10, length,30);
		stage.addChild(rect);

		var text = new createjs.Text(message, "13px Mono", "black");
		text.x = x+15;
		text.y = y+22;
		text.textBaseline = "alphabetic";
		stage.addChild(text);
		stage.update();
		//to remove
		currTooltip=rect;
		currTooltipText=text;
	}

	function createCommit(stage,x,y,message){
		var circle = new createjs.Shape();
		circle.graphics.setStrokeStyle(1).beginStroke("Black").beginFill("DeepSkyBlue").drawCircle(0, 0, 6);
		circle.x = x;
		circle.y = y;
		circle.addEventListener('mouseover',function (event){
			circle.graphics.clear().setStrokeStyle(1).beginStroke("Black").beginFill("Red").drawCircle(0, 0, 6);
			stage.update();
			createTooltip(stage,x,y,message);
		});

		circle.addEventListener('mouseout',function (event){
			circle.graphics.clear().setStrokeStyle(1).beginStroke("Black").beginFill("DeepSkyBlue").drawCircle(0, 0, 6);
			stage.removeChild(currTooltip);
			stage.removeChild(currTooltipText);
			stage.update();
		});

		stage.addChild(circle);
		return {
			position:{
				x:x,
				y:y
			},
			message:message
		}
	}

	function connectCommit(stage,cmtSrc, cmtDst){
		var line = new createjs.Shape();
		line.graphics.setStrokeStyle(2);
		if (cmtSrc.position.y == cmtDst.position.y){
			var width = cmtDst.position.x-cmtSrc.position.x;
			line.graphics.beginStroke("Black");
			line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
			line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
		}else{
			var deltaY = cmtDst.position.y-cmtSrc.position.y;
			var deltaX = cmtDst.position.x-cmtSrc.position.x;				
			if (cmtSrc.position.y < cmtDst.position.y){
				line.graphics.beginStroke("Blue");
				line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x,cmtDst.position.y-15);
				line.graphics.lineTo(cmtSrc.position.x+deltaX-10,cmtDst.position.y-15);
				line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
			}
			if (cmtSrc.position.y > cmtDst.position.y){
				line.graphics.beginStroke("Red");
				line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x+deltaX/2,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x+deltaX/2,cmtSrc.position.y+deltaY/2);
				line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
			}
		}
		stage.addChildAt(line,0);
	}

	function drawDate(stage,y,x1,x2,color,date){
		var rect = new createjs.Shape();
		rect.graphics.setStrokeStyle(1).beginStroke("White").beginFill(color).drawRect(x1, y, x2-x1,20);
		stage.addChild(rect);
		var text = new createjs.Text(date, "12px Mono", "white");
		text.x = x1+((x2-x1)/2)-7;
		text.y = y+14;
		text.textBaseline = "alphabetic";
		stage.addChild(text);
	}
