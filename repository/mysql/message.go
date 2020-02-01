//Package message ...
package mysql

import (
	"database/sql"
	"github.com/amiraliio/tgbp-api/config"
	"github.com/amiraliio/tgbp-api/domain/channel"
	"github.com/amiraliio/tgbp-api/domain/message"
	"github.com/amiraliio/tgbp-api/domain/user"
	"github.com/pkg/errors"
)

//TODO pagination for getAll messages

type messageRepo struct {
	db *sql.DB
	//TODO mysql client
	//TODO mysql close
	//TODO APP Config
}

func NewMysqlMessageRepository(db *sql.DB) message.MessageRepository {
	return &messageRepo{
		db,
	}
}

func (m *messageRepo) DirectMessagesList(userID, receiverID, channelID int64) ([]*message.Message, error) {
	app := new(config.App)
	app = app.SetAppConfig()
	db := app.DB()
	defer db.Close()
	rows, err := db.Query("select me.userID, me.message, me.createdAt, uu.username,cha.channelName, cha.channelType from messages as me inner join users as us on me.userID=us.userID inner join users_usernames as uu on uu.userID=us.id and me.channelID=uu.channelID inner join channels as cha on cha.id=me.channelID where me.type=? and me.channelID=? and ((me.userID=? and me.receiver=?) or (me.receiver=? and me.userID=?)) order by me.createdAt asc, me.id asc", "DM", channelID, userID, receiverID, userID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*message.Message
	for rows.Next() {
		messageModel := new(message.Message)
		userModel := new(user.User)
		username := new(user.UserUserName)
		channelModel := new(channel.Channel)
		if err := rows.Scan(&messageModel.UserID, &messageModel.Message, &messageModel.CreatedAt, &username.Username, &channelModel.ChannelName, &channelModel.ChannelType); err != nil {
			return nil, errors.Wrap(err, "repository.mysql.Message.DirectMessagesList")
		}
		if messageModel.UserID == userID {
			username.Username = "You"
		} else {
			username.Username = "[User " + username.Username + "]"
		}
		userModel.UserSign = username
		messageModel.User = userModel
		messageModel.Channel = channelModel
		messages = append(messages, messageModel)
	}
	return messages, nil
}
