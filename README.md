### Tournament
It is a tournament service. Each player holds certain amount of bonus points, which can be spent for goods or for
joining tournament. Player can join only if they have enough money for pay tournament deposit.

The service has 5 endpoints:
1. Take and fund player: /take?playerId=1&points=300 take 300 points from player; /fund?playerId=1&points=300 
funds player 1 with 300 points.
2. Announce tournament specifying the entry deposit: /announceTournament?tournamentId=1&deposit=1000
3. Join player into a tournament: /joinTournament?tournamentId=1&playerId=1. A player play on his own money.
4. Result tournament winners and prizes: /resultTournament?tournamentId=1, 
  response: {"winners":[{"playerId":"1","prize":500,"balance":600}]}
5. Player balance: /balance?playerId=1, response: {"playerId":"1", "points":"500"}

If player does not exist, fund endpoint create them with balance=points. After tournament results winner is choosen
 randomly and gets prize.
Endpoints 1-4 return HTTP status codes only like 2xx, 4xx, 5xx (when /fund create new player, it also returns json
format of them). Endpoint 5 returns json format of winners.

That service has wroten package postgres for working with database. If you use it, you will need to create two tables:
1. tournaments, which has following columns: id text primary key, deposit integer > 0, prize integer >= 0, participants
 text array, winner json, isOpen bool (shows tournament current state)
2. players with following columns: id text primary key, points integer >= 0
