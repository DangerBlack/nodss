# NoDSS
## Node DSS (Decision Support System)

Is a completely new implementation of [https://github.com/DangerBlack/Google-Decision-Support-System](https://github.com/DangerBlack/Google-Decision-Support-System) with a more permissive licence (M.I.T.)

## Telegram Bot

This code is the backend of a bot named [@dsspiebot](http://telegram.me/dsspiebot).
In order to run the bot locally you must create a file named `.env`

With the Token of the bot.

```
TOKEN=xxxxxxxxxx:zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
```

### Start node
To start the bot just write

```
npm i
npm run start
```

### Docker container
To start the docker container

```
docker build -t nodss .
docker run -d nodss
```

## Library

The library can be imported in other project to perform some google query. 
Pay attention google may detect you are doing something against some deep google policy so do not abuse it.

## Why?

It's really useful to understand the global mood about some important event in the world.
Can be used to predict political outcome based on what people wrote on the internet, but also to decide a new hair cut to impress a girl.
The most loved name of a Dog or if you should or should not invest in Tesla...