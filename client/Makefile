CXX = g++
CXXFLAGS = -std=c++17 -Wall -I./websocketpp -I/usr/include/jsoncpp -I/usr/include -I/usr/include/boost
LDFLAGS = -ljsoncpp -lncurses -lpthread -lboost_system

all: tinytetris tinytetris-commented

tinytetris-commented: tinytetris-commented.cpp
		$(CXX) $(CXXFLAGS) tinytetris-commented.cpp -o tinytetris-commented $(LDFLAGS)

tinytetris: tinytetris.cpp
		$(CXX) $(CXXFLAGS) tinytetris.cpp -o tinytetris $(LDFLAGS)

clean:
		rm -f tinytetris tinytetris-commented