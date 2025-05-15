# Floury Backend
## 🚧 This project is still a heavy WIP 🚧

This is the backend of my website. It's mainly used for contact handling (basic SMTP) and in the future handling of blog posts and commenting on the website. 

### Usage
You can use this docker-compose to run it, or just build from source and run on bare metal
```
version: "3.8"

services:
  backend:
    image: ghcr.io/monkakokosowa/backend-flour:main
    container_name: backend-flour
    env_file: .env
    ports:
      - "4324:${PORT}"
    restart: unless-stopped
```
Remember to rename .env-example file and fill it with your own values!


### Issues
If you find any bugs, please don't abuse them but just make a bug report at [the issues tab](https://github.com/MonkaKokosowa/backend-flour/issues). 
There's not much to abuse anyways, it's just me here :p

