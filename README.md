# ping-webhook-bot

@ping_webhook_bot on telegram. Converts your get requests into telegram messages with formatting and timestamps.

## Why

Might be useful if you need are running some technology that is new to you or you don't have time/resources to introduces proper notifications for something.

Like when you have legacy form ops with parts of monitoring / aggregation done by bash scripts and you want realtime notifications without doing more than adding 
```bash
nodes=114
curl "https://wh.v-trof.ru/MGTIiN1mR?m=nodesOnline=$nodes&t=info"
```
to those scripts.


Same stuff with reporting erros of that 1 live example on your docs page, a simple fetch might be more convinient than setting up sentry for this.

## Usage

Mesage structure
`wh.v-trof.ru/{project-id}?m={message-to-send}&t={message-type}`

Thre are special symbols for some values of `message-type`.
- done — ✅
- warn — ⚠️
- info — ℹ️
- err — ⁉

Make sure to repalce space (` `) with `%20` as you are just making a get request after all.