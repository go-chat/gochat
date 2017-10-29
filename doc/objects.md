# User attributes

| field     | type   | description                                  |
|-----------|--------|----------------------------------------------|
| id        | string | id of user, this id is unique in application |
| name      | string | name of user                                 |
| image_url | string | avatar of user                               |
| meta      | json   | meta data for each application               |

# Group attributes

| field     | type         | description                                           |
|-----------|--------------|-------------------------------------------------------|
| id        | string       | id of a group chat, this is unique id for application |
| name      | string       | name a a group                                        |
| image_url | string       | avatar of a group                                     |
| meta      | json         | meta data for a group.                                |
| members   | list of user | list of users who is member of this group             |

# Message attributes

| field     | type   | description                                               |
|-----------|--------|-----------------------------------------------------------|
| id        | string | id of a message                                           |
| group_id  | string | id of group this message belongs to                       |
| sender_id | string | id of sender of this message                              |
| type      | int    | type of message                                           |
| content   | json   | content of message, this is depend of type of the message |

Text message: `type` = 0, content is a dictionary with these fields

| field | type   | description            |
|-------|--------|------------------------|
| text  | string | content of the message |

Image message: `type` = 1, content is a dictionary with these fields

| field     | type   | description  |
|-----------|--------|--------------|
| image_url | string | url of image |

Custom message: `type` = 2, custom type. Dictionary is a json. 
Custom type allow application could specified new type of message. 
E.g: send file, send sticker, bot message, new user join group chat, ...

# Event attributes

| field     | type       | description                                                 |
|-----------|------------|-------------------------------------------------------------|
| sender_id | string     | id of user who trigger this event. this value could be null |
| group_id  | string     | id of group this event belong to                            |
| type      | integer    | type of event                                               |
| content   | dictionary | content of event                                            |

The different between message and event is: message is persistent, but event is not.
Event is used for delivery real time event: new message to group, typing indicator, delivery indicator, seen indicator, ...