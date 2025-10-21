CC = x86_64-w64-mingw32-gcc
SRC = src/main.c src/utils.c
OUT = yafl.exe
OBJ = build/yafl.o

build:
	$(CC) -std=c99 -o $(OUT) $(SRC) -Wall -Werror
