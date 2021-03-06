package main

import (
	"database/sql"
	"github.com/oliverpauffley/chess_ladder/laddermethods"
	"github.com/oliverpauffley/chess_ladder/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type Mockdb struct {
	users       map[string]models.CredentialsInternal
	ladders     map[int]models.Ladder
	ladderUsers map[int]models.LadderUser
	models.Datastore
}

func (db *Mockdb) CreateUser(username, email, password string) error {
	var mockcredentials models.CredentialsInternal
	mockcredentials.Username = username
	mockcredentials.Hash, _ = bcrypt.GenerateFromPassword([]byte(password), 8)
	mockcredentials.Email = email
	db.users[username] = mockcredentials
	return nil
}

func (db *Mockdb) QueryByEmail(email string) (models.CredentialsInternal, error) {
	for _, entry := range db.users {
		if entry.Email == email {
			user := db.users[entry.Username]
			userCredentials := models.CredentialsInternal{Id: user.Id, Username: user.Username, Email: user.Email, JoinDate: user.JoinDate.Round(time.Hour),
				Role: user.Role, Wins: user.Wins, Losses: user.Losses, Draws: user.Draws, Hash: user.Hash}
			return userCredentials, nil
		}
	}
	return models.CredentialsInternal{}, sql.ErrNoRows
}

func (db *Mockdb) QueryById(id int) (models.CredentialsExternal, error) {
	for _, entry := range db.users {
		if entry.Id == id {
			user := db.users[entry.Username]
			userCredentials := models.CredentialsExternal{Id: user.Id, Username: user.Username, Email: user.Email, JoinDate: user.JoinDate.Round(time.Hour),
				Role: user.Role, Wins: user.Wins, Losses: user.Losses, Draws: user.Draws}
			return userCredentials, nil
		}
	}
	return models.CredentialsExternal{}, sql.ErrNoRows
}

func (db *Mockdb) DeleteUser(id int) error {
	user, err := db.QueryById(id)
	if err != nil {
		return err
	}
	delete(db.users, user.Username)
	return nil
}

func (db *Mockdb) AddLadder(name, method string, owner int) error {
	// find max key value
	key := 0
	for range db.ladders {
		if _, exists := db.ladders[key+1]; exists == true {
			key++
		}
	}
	// add new ladder
	newLadder := models.Ladder{Id: key, Name: name, Owner: owner, Method: method, HashId: "Ladder1"}
	db.ladders[key] = newLadder
	return nil
}

func (db *Mockdb) GetLadderFromHashId(hashId string) (models.Ladder, error) {
	id, err := strconv.Atoi(hashId)
	if err != nil {
		return models.Ladder{}, err
	}
	ladder, exists := db.ladders[id]
	if !exists {
		return models.Ladder{}, sql.ErrNoRows
	}
	return ladder, nil
}

func (db *Mockdb) JoinLadder(ladderId, userId int, method laddermethods.LadderMethod) error {
	// find key vale
	key := 0
	for range db.ladderUsers {
		if _, exists := db.ladders[key+1]; exists == true {
			key++
		}
	}
	// add user to ladder
	ladderUser := models.LadderUser{Id: key, LadderId: ladderId, UserId: userId, Rank: 0, Points: 1000}
	db.ladderUsers[key] = ladderUser
	return nil
}

func (db *Mockdb) GetLadders(userId int) ([]models.LadderInfo, error) {
	var userLadders []models.LadderInfo
	var laddersOwned []models.Ladder
	for _, ladder := range db.ladders {
		if ladder.Owner == userId {
			laddersOwned = append(laddersOwned, ladder)
		}
	}
	for _, ladder := range laddersOwned {

		var playerList []models.LadderRanks
		for _, player := range db.ladderUsers {
			if player.LadderId == ladder.Id {
				info, _ := db.QueryById(player.UserId)
				playerInfo := models.LadderRanks{
					Name:   info.Username,
					UserId: player.UserId,
					Rank:   player.Rank,
					Points: player.Points,
				}
				playerList = append(playerList, playerInfo)
			}
		}
		ladderWithPlayers := models.LadderInfo{
			LadderId: ladder.Id,
			Name:     ladder.Name,
			Owner:    ladder.Owner,
			HashId:   ladder.HashId,
			Players:  playerList,
		}
		userLadders = append(userLadders, ladderWithPlayers)
	}
	return userLadders, nil
}
