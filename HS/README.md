## GameServer
Implementation of a card game for a multi-user game server.

### Game Instructions 

//TODO


### Testing 

To build and execute the binary locally run:

```
make up
```
To spin up a docker container and execute the binary run:

```
make deploy
```
To see the game server's API in action, go to the tests/scripts/ directory and play with the scripts.
A typical work flow for manual tests would look something like this:

```
22:36 $ ./getRooms.sh
{
  "rooms": [
    {
      "room_id": "9589f7ff-0ae0-49c4-8460-8d24a4c16c1b",
      "players": 3
    },
    {
      "room_id": "b0dd44a2-619c-451a-acae-76605254326a",
      "players": 2
    },
    {
      "room_id": "f15f80ac-825b-43e1-aeff-c7110d534fbb",
      "players": 4
    },
    {
      "room_id": "d33b4d83-68d6-4823-bcf8-0744bf2f6d98",
      "players": 1
    },
    {
      "room_id": "cb11c3fc-9d4a-4954-be9e-afb5858ee1bd",
      "players": 2
    }
  ]
}
✔ ~/go/src/github.com/adrianosela/GameServer/tests/scripts [master|✔]
22:36 $ ./getPlayers.sh
Enter the Room ID: 9589f7ff-0ae0-49c4-8460-8d24a4c16c1b
{
  "players": [
    "4428913d-852a-4310-8fe4-835e374df07a",
    "f8f4d29b-7a2c-4020-a7fb-57244bf9b08d",
    "8af50a82-4fc1-4b74-b1bf-c0b45cc61b04"
  ]
}
✔ ~/go/src/github.com/adrianosela/GameServer/tests/scripts [master|✔]
22:36 $ ./getHand.sh
Enter the Room ID: 9589f7ff-0ae0-49c4-8460-8d24a4c16c1b
Enter the User ID: f8f4d29b-7a2c-4020-a7fb-57244bf9b08d
{
  "cards": [
    {
      "Type": "Queen",
      "Suit": "Diamond",
      "FaceUp": false,
      "VisibleToOwner": true
    },
    {
      "Type": "Ace",
      "Suit": "Club",
      "FaceUp": false,
      "VisibleToOwner": true
    },
    {
      "Type": "Six",
      "Suit": "Club",
      "FaceUp": false,
      "VisibleToOwner": true
    },
    {
      "Type": "Seven",
      "Suit": "Diamond",
      "FaceUp": false,
      "VisibleToOwner": false
    },
    {
      "Type": "Five",
      "Suit": "Diamond",
      "FaceUp": false,
      "VisibleToOwner": false
    },
    {
      "Type": "Jack",
      "Suit": "Diamond",
      "FaceUp": false,
      "VisibleToOwner": false
    }
  ]
}
```