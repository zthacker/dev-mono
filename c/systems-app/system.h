#ifndef SYSTEM_H
#define SYSTEM_H

class System {
public:
    virtual ~System() = default;
    virtual void update(double deltaTime) = 0;
    
};

#endif //SYSTEM_H