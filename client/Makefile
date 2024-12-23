CXX = g++
CXXFLAGS = -std=c++17 -Wall -I/usr/include/jsoncpp -I/usr/include
LDFLAGS = -ljsoncpp -lncurses -lpthread -lboost_system

all: tinytetris tinytetris-commented

tinytetris: tinytetris.cpp
	$(CXX) $(CXXFLAGS) tinytetris.cpp -o tinytetris $(LDFLAGS)

tinytetris-commented: tinytetris-commented.cpp
	$(CXX) $(CXXFLAGS) tinytetris-commented.cpp -o tinytetris-commented $(LDFLAGS)

clean:
	rm -f tinytetris tinytetris-commented
