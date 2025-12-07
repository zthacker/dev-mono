#ifndef ACTION_SYSTEM_H
#define ACTION_SYSTEM_H

#include "system.h"

class ActionSystem : public System {
public:
    void update(double deltaTime) override;
};

#endif //ACTION_SYSTEM_H