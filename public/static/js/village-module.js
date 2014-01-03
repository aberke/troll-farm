
/* -------------------------- VILLAGE ---------------------- */

var TrollConnection = function() {
	var ws;


	function connect() {
		var ws = new WebSocket(WEBSOCKET_HOST + NEW_CONNECTION_ENDPOINT);
		return ws;
	}
	this.init = function(messageCallbacks, callback) {
		ws = connect();

		ws.onmessage = function(event) {
			var msg = event.data;
			if (typeof msg === "string"){ msg = JSON.parse(msg); }

			console.log(msg["Type"])

			if (messageCallbacks[msg.Type])
				messageCallbacks[msg.Type](msg);
			else
				console.log("Recieved unrecognized message type: " + msg.Type);
		};
		ws.onopen = function(event) {
			if (callback) callback();
		}
	}
	function send(msgBody) {
		ws.send(JSON.stringify(msgBody));
	}
	this.sendPing = function() {
		var msg = {"message-type": "ping",};
		send(msg);
	}
	this.sendTrollsRequest = function() {
		var msg = {"message-type": "trolls"};
		send(msg);
	}
	this.sendMove = function(x, y) {
		var msg = { "message-type": "move",
					"data": {"x":String(x), "y":String(y)}
				  }
		send(msg);
	}
}

var TrollVillageModule = function(widgetDiv) {
	var self = this;

	this.items = {};  // maps {itemID: DrawableItem}
	this.localID;

	var trollConnection;

	this.widgetDiv = widgetDiv;
	this.canvas;
	this.movingContext;
	this.staticContext;

	/* board is a 10x10 grid with 40x40 px squares  */
	this.board = {"width": 10,
				  "height": 10,
				  "cellSize": 40,
				 }

	this.removeItem = function(trollID) {
		if (this.items[trollID]) {
			this.items[trollID].erase();
			delete this.items[trollID];
		}
	}
	this.updateItem = function(itemID, item) {
		if (item.Name == "DELETE") 
			return this.removeItem(itemID);

		if (self.items[itemID]) {
			self.items[itemID].update(item.Coordinates.x, item.Coordinates.y);
		} else {
			self.items[itemID] = new OtherTroll();
			self.items[itemID].init(item.Coordinates.x, item.Coordinates.y);
		}	
	}

	this.recieveUpdate = function(msg) {
		console.log('TrollVillageModule recieveUpdate')
		console.log(msg)

		var item; // recycled variable as iterate through map
		var itemsMap = msg.ItemsMap;

		for (var itemID in itemsMap) {
			item = itemsMap[itemID]
			self.updateItem(itemID, item);
		}	
	}

	this.recieveItems = function(msg) {
		console.log('TrollVillageModule recieveItems')
		console.log(msg)

		// clear out old items
		self.items = {};

		self.localID = msg.LocalTroll;
		var item; // recycled variable as iterate through map

		var itemsMap = msg.ItemsMap;
		for (var itemID in itemsMap) {

			var newItem;
			// anything that doesnt have a special item name is a Troll
			item = itemsMap[itemID];
			if (item.Name == "BANANA") {
				newItem = new Banana();
			} else if (item.Name == "FOODBUTTON") {
				newItem = new FoodButton();
			} else if (itemID == self.localID) {
				newItem = new LocalTroll();
			} else {
				newItem = new OtherTroll();
			}
			newItem.init(item.Coordinates.x, item.Coordinates.y);
			newItem.id = itemID;
			self.items[itemID] = newItem;

		}
	}
	this.recievePing = function(msg) {
		console.log("ping -> pong");
	}

	this.createCanvas = function() {
		var canvas = document.createElement('canvas');
		canvas.width = this.board.width*this.board.cellSize;
		canvas.height = this.board.height*this.board.cellSize;
		return canvas;
	}
	this.drawBoard = function() {
		this.staticCanvas = this.createCanvas();
		this.staticContext = this.staticCanvas.getContext("2d");
		this.widgetDiv.appendChild(this.staticCanvas);

		this.movingCanvas = this.createCanvas();
		this.movingContext = this.movingCanvas.getContext("2d");
		this.widgetDiv.appendChild(this.movingCanvas);

		for (var x=0; x<=this.board.width; x+=1) {
			this.staticContext.moveTo(x*this.board.cellSize, 0);
			this.staticContext.lineTo(x*this.board.cellSize, this.board.height*this.board.cellSize);
		}
		for (var y=0; y<=this.board.height; y+=1) {
			this.staticContext.moveTo(0, y*this.board.cellSize);
			this.staticContext.lineTo(this.board.width*this.board.cellSize, y*this.board.cellSize);
		}	
		this.staticContext.strokeStyle = "black";
		this.staticContext.stroke();
	}

	/* Define an object to hold all our images for the game so images are only ever created once. */
	var imageRepository = new function() {
        // Define images
        this.troll 		= new Image();
        this.otherTroll = new Image();
        this.foodButton = new Image();
        this.banana		= new Image();

        // Ensure all images have loaded before starting the game
        var numImages = 4;
        var numLoaded = 0;
        function imageLoaded() {
            numLoaded++;
            if (numLoaded === numImages) {
                    self.init();
            }
        }
        this.troll.onload = function() {
                imageLoaded();
        }
        this.otherTroll.onload = function() {
                imageLoaded();
        }
        this.foodButton.onload = function() {
                imageLoaded();
        }
        this.banana.onload = function() {
                imageLoaded();
        }

        // Set images src
        this.troll.src 		= DOMAIN + "/static/img/troll.gif";
        this.otherTroll.src = DOMAIN + "/static/img/other-troll.gif";
        this.foodButton.src = DOMAIN + "/static/img/troll-food-button.JPG";
        this.banana.src 	= DOMAIN + "/static/img/banana.gif";
	}
	this.init = function() {
		/* init is called only after the imageRepository has loaded in all images */
		this.drawBoard();
		Drawable.prototype.board   		= this.board;
		Drawable.prototype.context 		= this.movingContext;
		StaticDrawable.prototype.context= this.staticContext;
		LocalTroll.prototype.img 		= imageRepository.troll;
		OtherTroll.prototype.img 		= imageRepository.otherTroll;
		Banana.prototype.img 	 		= imageRepository.food;
		FoodButton.prototype.img 		= imageRepository.foodButton;


		trollConnection = new TrollConnection();
		trollConnection.init({"items": this.recieveItems,
							  "update": this.recieveUpdate,
							  "ping": this.recievePing} );
	}

function Drawable() {
	this.x;
	this.y;
	this.padding = 5;
	this.x_px;
	this.y_px;
	this.width = 20;
	this.height = 20;
	this.board;

	this.img;
	this.context;

	this.erase = function() {
		this.context.clearRect(this.x_px, this.y_px, this.width, this.height);
	}

	this.draw = function() {
		console.log(this)
		if (this.x_px)
			this.erase();

		this.x_px = this.x*this.board.cellSize + this.padding;
		this.y_px = this.y*this.board.cellSize + this.padding;
		
		this.context.drawImage(this.img, this.x_px, this.y_px, this.width, this.height);
	}
	this.setImage = function(img) {
		this.img = img;
	}
	this.update = function(x, y) {
		this.x = x;
		this.y = y;
		this.draw();
	}
	this.init = function(x, y) {
		this.x = x;
		this.y = y;
		this.draw();
	}
}
function StaticDrawable() {
}
StaticDrawable.prototype = new Drawable();
function Banana() {
}
Banana.prototype = new Drawable();
function FoodButton() {
	this.padding = 1;
	this.width = this.board.cellSize - 2;
	this.height = this.board.cellSize - 2;
}
FoodButton.prototype = new StaticDrawable();

function Troll() {

	this.print = function() {
		console.log(this)
	}

	this.move = function(direction) {
		if (direction == "left") {
			trollConnection.sendMove(-1, 0);
		} else if (direction == "right") {
			trollConnection.sendMove(1, 0);
		} else if (direction == "up") {
			trollConnection.sendMove(0, -1);
		} else if (direction == "down") {
			trollConnection.sendMove(0, 1);
		} else {
			console.log("direction: " + direction);
		}
	}
}
Troll.prototype = new Drawable();
function LocalTroll() {
	var self = this;

	// need closure
	var onkeydown = function(key) {
		self.move(key);
	}
	setKeyListeners(onkeydown);
}
LocalTroll.prototype = new Troll();

function OtherTroll() {
}
OtherTroll.prototype = new Troll();






}


/* --------------------------------------------------------------- */

KEY_CODES = {
  32: 'space',
  37: 'left',
  38: 'up',
  39: 'right',
  40: 'down',
}
function setKeyListeners(onkeydownCall) {

	document.onkeydown = function(e) {
		// Firefox and opera use charCode instead of keyCode to return which key was pressed.
		var keyCode = (e.keyCode) ? e.keyCode : e.charCode;
		if (KEY_CODES[keyCode]) {
			e.preventDefault();
			onkeydownCall(KEY_CODES[keyCode]);
		}
	}
	document.onkeyup = function(e) {
	  var keyCode = (e.keyCode) ? e.keyCode : e.charCode;
	  if (KEY_CODES[keyCode]) {
	    e.preventDefault();
	  }
	}


}


