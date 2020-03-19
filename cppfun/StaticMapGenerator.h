//
// Created by selva on 3/14/20.
//

#ifndef CPPFUN_STATICMAPGENERATOR_H
#define CPPFUN_STATICMAPGENERATOR_H

#include <string>
#include <map>
#include <utility>
#include <initializer_list>

enum class Fruits {
    Mango,
    Apple,
    Jamun,
    Grapes
};

enum class Season {
    Summer,
    Winter,
    Rainy,
    Spring,
    AllYear
};

#define MAKE_MAP(map_name, key_type, value_type, default_value) \
template <key_type, value_type  default_val = default_value> \
struct map_name { \
    static const value_type val = default_val; \
};

#define ADD_KEY_VALUE(map_name,key,value) \
template <> \
struct map_name<key> { \
    static const  decltype(value) val = value; \
};

MAKE_MAP(FruitMap,Fruits,Season,Season::AllYear);
ADD_KEY_VALUE(FruitMap,Fruits::Mango,Season::Summer);

template <typename key_type, typename value_type,value_type default_value>
class DefaultValueMap {
public:
    DefaultValueMap(std::initializer_list<std::pair<const key_type,value_type>> key_value_pair ):
    KeyValueMap{key_value_pair}
    {}
    value_type getValue(const key_type& key) const noexcept {
        auto itr = KeyValueMap.find(key);
        if(itr == KeyValueMap.end()){
            return default_value;
        }
        return itr->second;
    }
private:
    std::map<key_type,value_type> const KeyValueMap;
};
#endif //CPPFUN_STATICMAPGENERATOR_H
