build:
	x86_64-w64-mingw32-gcc -std=c99 -o main src/main.c src/utils.c -Wall -Werror