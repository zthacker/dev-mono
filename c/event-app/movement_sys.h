//
// Created by Zach Thacker on 12/5/25.
//

#ifndef MOVEMENT_SYS_H
#define MOVEMENT_SYS_H
#include <iostream>
#include <map>
#include <ostream>
#include <string>

struct Vector {
    int x;
    int y;
    int z;
};

struct Position {
    std::string name;
    Vector position;
};

class MovementSystem {
    public:
    MovementSystem();
    ~MovementSystem();

     void onMove(const Position& t) {
        std::cout << t.name << "moved to "
        << "x: " << t.position.x
        << " y: " << t.position.y
        << " z: " << t.position.z
        << std::endl;

        if (positions.find(t.name) == positions.end()) {
            positions[t.name] = t.position;
        } else {
            positions[t.name] = t.position;
        }
    }

private:
    //entities
    std::map<std::string, Vector> positions;
};

#endif //MOVEMENT_SYS_H
