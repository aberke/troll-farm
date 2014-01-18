Troll Farm
====================
<img src="https://troll-farm.herokuapp.com/static/img/other-troll.gif"
 alt="troll-farm logo" title="A real troll" align="right" />
<img src="https://troll-farm.herokuapp.com/static/img/other-troll.gif"
 alt="troll-farm logo" title="A real troll" align="right" />

By trolls, for trolls.

<https://troll-farm.herokuapp.com/>

Why
----
I want a pleasant space, isolated environment for trolls to roam freely, meet each other, and communicate.  They should have the opportunity to spawn and aspire to grow long beards.

I mean I wanted to learn Go and play with web sockets and make something interesting.


TODO
---

* Allow trolls to move from Grid to Grid

* TODO: MAKE WORK ON MOBILE

* Test sending bad messages to server.
AKA test error handling of malformed/malicious 

* Write tests


Structure
---

Each troll-client (Troll) handles its own websocket connection

The server has (among other things like channels): 

* a map of trolls [trollID -> *Troll]
* a map of grids  [gridID  -> *Grid]	
* a map from troll to grid it lives in [trollID -> gridID]

Grid -> Troll is a 1-to-many relationship.  Many Trolls and other items (see GRID_CAPACITY) live in a Grid

update messages should only be sent to the trolls in the grid that recieved an update















