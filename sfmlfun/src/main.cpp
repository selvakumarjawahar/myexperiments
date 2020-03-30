//
// Created by selva on 21-03-2020.
//
#include <iostream>

#include "SFML/Graphics.hpp"
#include "SFML/System.hpp"
#include <fstream>

sf::Vector2f viewSize(1024, 768);
sf::VideoMode vm(viewSize.x, viewSize.y);
sf::RenderWindow window(vm, "Hello SFMLGame !!!", sf::Style::Default);


enum class Directions {
  Up,
  Down,
  Right,
  Left,
  NoMove
};
bool playerMoving = false;
Directions direction = Directions::NoMove;

void updateInput()
{
  sf::Event event;
  while (window.pollEvent(event)) {
    if (event.type == sf::Event::KeyPressed) {
      playerMoving = true;
      direction = Directions::NoMove;
      if (event.key.code == sf::Keyboard::Right) {
        direction = Directions::Right;
      }
      if (event.key.code == sf::Keyboard::Left) {
        direction = Directions::Left;
      }
      if (event.key.code == sf::Keyboard::Up) {
        direction = Directions::Up;
      }
      if (event.key.code == sf::Keyboard::Down) {
        direction = Directions::Down;
      }
    }

    if (event.type == sf::Event::KeyReleased) {
      playerMoving = false;
    }
    if (event.key.code == sf::Keyboard::Escape || event.type == sf::Event::Closed)
      window.close();
  }
}

sf::Vector2f getCoordinates(float dt, Directions dir)
{
  auto displacement = 200.0f * dt;
  switch (dir) {
  case Directions::Right:
    return sf::Vector2f(displacement, 0);
  case Directions::Left:
    return sf::Vector2f(-displacement, 0);
  case Directions::Up:
    //return sf::Vector2f(0, -displacement);
  case Directions::Down:
    //return sf::Vector2f(0, displacement);
  default:
    return sf::Vector2f(0, 0);
  }
}

int main()
{

  sf::RectangleShape rect(sf::Vector2f(500.0f, 300.0f));
  rect.setFillColor(sf::Color::Yellow);
  rect.setPosition(viewSize.x / 2, viewSize.y / 2);
  rect.setOrigin(sf::Vector2f(rect.getSize().x / 2, rect.getSize().y / 2));

  sf::Texture sky_texture;
  sf::Sprite sky_sprite;
  sf::Image sky_image;

  sf::Texture bg_texture;
  sf::Sprite bg_sprite;
  sf::Image bg_image;

  sf::Texture hero_texture;
  sf::Sprite hero_sprite;
  sf::Image hero_image;

  sf::Texture golem_texture;
  sf::Sprite golem_sprite;
  sf::Image golem_image;

  sf::Texture bullet_texture;
  sf::Sprite bullet_sprite;
  sf::Image bullet_image;

  if (!sky_image.loadFromFile("C:/Users/selva/Projects/myexperiments/sfmlfun/src/assets/graphics/sky.png")) {
    std::cout << "Error in loading sky image " << '\n';
    return 0;
  }
  if (!sky_texture.loadFromImage(sky_image)) {
    std::cout << "Error in loading sky texture " << '\n';
    return 0;
  }
  sky_sprite.setTexture(sky_texture);

  if (!bg_image.loadFromFile("C:/Users/selva/Projects/myexperiments/sfmlfun/src/assets/graphics/bg.png")) {
    std::cout << "Error in loading bg image " << '\n';
    return 0;
  }
  if (!bg_texture.loadFromImage(bg_image)) {
    std::cout << "Error in loading bg texture " << '\n';
    return 0;
  }
  bg_sprite.setTexture(bg_texture);

  if (!hero_image.loadFromFile("C:/Users/selva/Projects/myexperiments/sfmlfun/src/assets/graphics/heroAnim.png")) {
    std::cout << "Error in loading hero image " << '\n';
    return 0;
  }
  if (!hero_texture.loadFromImage(hero_image)) {
    std::cout << "Error in loading hero texture " << '\n';
    return 0;
  }
  hero_sprite.setTexture(hero_texture);
  hero_sprite.setTextureRect(sf::IntRect(0,0,92,126));
  hero_sprite.setPosition(sf::Vector2f(viewSize.x / 2, 500));
  hero_sprite.setOrigin(hero_texture.getSize().x / 2, hero_texture.getSize().y / 2);

  if (!golem_image.loadFromFile("C:/Users/selva/Projects/myexperiments/sfmlfun/src/assets/graphics/golem.png")) {
    std::cout << "Error in loading golem image " << '\n';
    return 0;
  }
  if (!golem_texture.loadFromImage(golem_image)) {
    std::cout << "Error in loading golem texture " << '\n';
    return 0;
  }
  golem_sprite.setTexture(golem_texture);
  golem_sprite.setPosition(sf::Vector2f(viewSize.x/2,480));
  golem_sprite.setOrigin(golem_texture.getSize().x / 2, golem_texture.getSize().y / 2);

  if (!bullet_image.loadFromFile("C:/Users/selva/Projects/myexperiments/sfmlfun/src/assets/graphics/bullet.png")) {
    std::cout << "Error in loading bullet image " << '\n';
    return 0;
  }
  if (!bullet_texture.loadFromImage(bullet_image)) {
    std::cout << "Error in loading bullet texture " << '\n';
    return 0;
  }
  bullet_sprite.setTexture(bullet_texture);
  bullet_sprite.setPosition(sf::Vector2f(viewSize.x/2,200));
  bullet_sprite.setOrigin(bullet_texture.getSize().x / 2, bullet_texture.getSize().y / 2);

   sf::Clock clock;
   float elasped_time =0;
   const float anim_duration = 1.5f;
   const int framecount = 4;

  while (window.isOpen()) {
    // Handle Keyboard events
    // Update Game Objects in the scene
    updateInput();
    sf::Time delta = clock.restart();
    elasped_time += delta.asSeconds();
    int animate_frameno = (static_cast<int>((elasped_time/anim_duration)*framecount) % framecount);
    hero_sprite.setTextureRect(sf::IntRect(92*animate_frameno,0,92,126));
    if (playerMoving) {
      hero_sprite.move(getCoordinates(delta.asSeconds(), direction));
    }

    window.clear();
    // Render Game Objects
    window.draw(sky_sprite);
    window.draw(bg_sprite);
    window.draw(hero_sprite);
    window.draw(golem_sprite);
    window.draw(bullet_sprite);
    window.display();
  }
  return 0;
}
