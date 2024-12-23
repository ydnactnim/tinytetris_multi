# tinytetris

80x23 terminal tetris!

![tinytetris gif](animation.gif)

### tinytetris.cpp

This is the 80x23 version. You control it with `a` (left), `d` (right), `w` (rotate),
`s` (drop), and `q` (quit). It depends on `curses.h` (so you'll need to compile with
`-lcurses`, and install curses if you don't already have it) and requires C++11.

### tinytetris-commented.cpp

This one is almost identical to `tinytetris.cpp`, but not minified, and with some
comments to make it easier to read (but it's still tricky to read in certain parts).

### build binary tinytetris-commented.cpp

`g++ -o tinytetris-commented tinytetris-commented.cpp -lncurses`

tinytetris_multi/
├── server/ # Go로 작성된 서버 코드
│ ├── main.go # 서버의 진입점
│ ├── go.mod # Go 모듈 파일
│ └── ... # 기타 서버 관련 파일
├── client/ # 기존 C로 작성된 클라이언트 코드
│ ├── tinytetris.c # 클라이언트의 메인 파일
│ └── ... # 기타 클라이언트 관련 파일
├── README.md # 프로젝트 설명서
└── ...
