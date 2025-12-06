#include <iostream>
#include "bus.h"
#include "movement_sys.h"

int main() {
    std::cout << "Event Bus" << std::endl;

    Bus bus;
    MovementSystem movement_sys;

    // lambda for callback
    bus.subscribe<Position>([&](const Position& position) {
        movement_sys.onMove(position);
    });

    // publish events
    std::cout << "Event App - Movement System" << std::endl;

    bus.publish(Position{"User1", Vector{1,2,3}});
    bus.publish(Position{"User1", Vector{2,3,4}});
    bus.publish(Position{"User2", Vector{3,3,4}});
    bus.publish(Position{"User3", Vector{3,3,5}});
    bus.publish(Position{"User1", Vector{0,3,5}});


    return 0;
}
