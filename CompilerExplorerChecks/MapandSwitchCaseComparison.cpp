#include<iostream>
#include<map>

enum class Fruits {
    Mango,
    Grapes,
    Strawberry,
    Jamun,
    Watermelon,
    Banana,
    Kiwi,
    PineApple,
    Anjeer
};

enum class Season {
    Summer,
    Winter,
    Rainy,
    Spring,
    Autumn,
    Monsoon,
    Hot,
    AllSeason
};

static  const std::map<Fruits,Season> FruitSeasonMap = {
    {Fruits::Mango,Season::Summer},
    {Fruits::Grapes,Season::Spring},
    {Fruits::Strawberry,Season::Winter},
    {Fruits::Jamun,Season::Rainy},
    {Fruits::Watermelon,Season::Summer},
    {Fruits::Anjeer,Season::Monsoon},
    {Fruits::PineApple,Season::Hot},
    {Fruits::Kiwi,Season::Autumn}

};

const Season defaultval = Season::AllSeason;

Season getSeason(const Fruits fruit) {
    Season season;
    switch(fruit){
        case Fruits::Mango:
        season = Season::Summer;
        break;
        case Fruits::Grapes:
        season = Season::Spring;
        break;
        case Fruits::Strawberry:
        season = Season::Winter;
        break;
        case Fruits::Jamun:
        season = Season::Rainy;
        break;
        case Fruits::Watermelon:
        season = Season::Summer;
        break;
        case Fruits::PineApple:
        season = Season::Hot;
        break;
        case Fruits::Anjeer:
        season = Season::Monsoon;
        break;
        case Fruits::Kiwi:
        season = Season::Autumn;
        break;
        default:
        season = Season::AllSeason;
        break;
    }
    return season;
}

Season getSeasonFromMap(const Fruits fruit) {
    auto itr = FruitSeasonMap.find(fruit);  
    if(itr != FruitSeasonMap.end()){
        return itr->second;
    }
    return defaultval ;
}



 