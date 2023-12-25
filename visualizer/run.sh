#!/bin/bash

# Имя именованного канала
fifo="./visfifo"

# Создаем именованный канал
mkfifo "$fifo"