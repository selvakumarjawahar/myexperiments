#include <iostream>
#include "StaticMapGenerator.h"
#include <cstdlib>

Fruits FruitGenerator(){
    auto fruit = rand() % 3;
    return static_cast<Fruits>(fruit);
}
using DefaultFruitMap = DefaultValueMap<Fruits,Season,Season::AllYear>;

int main() {
    std::cout << "Default = " << (int) FruitMap<Fruits::Grapes>::val << '\n';
    std::cout << "Mango Season = " << (int) FruitMap<Fruits::Mango>::val << '\n';
    DefaultFruitMap dfmap{{Fruits::Mango,Season::Summer}};
    std::cout << "Mango Season = " << (int) dfmap.getValue(Fruits::Mango) << '\n';
    std::cout << "Random Fruit Season = " << (int) dfmap.getValue(FruitGenerator()) << '\n';

    return 0;
}
