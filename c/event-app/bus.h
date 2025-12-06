//
// Created by Zach Thacker on 12/5/25.
//

#ifndef BUS_H
#define BUS_H
#include <functional>
#include <map>
#include <typeindex>
#include <vector>
#include <memory>


class Bus {
    public:
    Bus();
    ~Bus();

    // register a callback for the EventType
    template<typename EventType>
    void subscribe(std::function<void(const EventType&)> callback);

    // dispatch event to all listeners
    template<typename EventType>
    void publish(const EventType& event);


    private:
    // handler base
    struct HandlerBase {
        virtual ~HandlerBase() = default;
    };

    // wrapper for callbacks for Type T
    template<typename T>
    struct HandlerList : HandlerBase {
        std::vector<std::function<void(const T&)>> callbacks;
    };

    // map of TypeID -> pointer to list
    std::map<std::type_index, std::unique_ptr<HandlerBase>> subscribers;

    // get the subscribers
    template<typename T>
    std::vector<std::function<void(const T&)>>& getSubscribers() {
        std::type_index typeIndex(typeid(T));

        if (subscribers.find(typeIndex) == subscribers.end()) {
            subscribers[typeIndex] = std::make_unique<HandlerList<T>>();
        }

        // cast
        return static_cast<HandlerList<T>*>(subscribers[typeIndex].get())->callbacks;
    }
};

// subscribe
template<typename EventType>
void Bus::subscribe(std::function<void(const EventType&)> callback) {
    auto& handlers = getSubscribers<EventType>();
    handlers.push_back(callback);
}

// publish
template<typename T>
void Bus::publish(const T& event) {
    auto& handlers = getSubscribers<T>();

    for (auto& handler : handlers) {
        handler(event);
    }
}



#endif //BUS_H
