#include <iostream>
#include "box2d/box2d.h"
#include "SFML/Graphics.hpp"
#include "SFML/System.hpp"

#define SCALE_FACTOR 10
#define WIDTH 1280.0f
#define HEIGHT 720.0f

//Below sizes in meters
#define BALL_RADIUS 1.0f
#define PADDLE_HEIGHT  1.0f
#define PADDLE_WIDTH 10.0f

sf::Vector2f viewSize(WIDTH, HEIGHT);
sf::VideoMode vm(viewSize.x,viewSize.y);
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

sf::Vector2f getForceVector(float dt, Directions dir)
{
  auto force = 0.1f;
  switch (dir) {
  case Directions::Right:
    return sf::Vector2f(force, 0);
  case Directions::Left:
    return sf::Vector2f(-force, 0);
  case Directions::Up:
    //return sf::Vector2f(0, -displacement);
  case Directions::Down:
    //return sf::Vector2f(0, displacement);
  default:
    return sf::Vector2f(0, 0);
  }
}

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

class MyContactListener : public b2ContactListener
{
public:
  /*
  void BeginContact(b2Contact* contact) {
    std::cout<< "Begin contact called" << '\n';
  }

  void EndContact(b2Contact* contact) {
    std::cout<< "End contact called" << '\n';
  }
   */

  void PreSolve(b2Contact* contact, const b2Manifold* oldManifold) {
    //contact->SetRestitution(1.0f);
    //std::cout<< "Presolve contact called" << '\n';
  }
  /*
  void PostSolve(b2Contact* contact, const b2ContactImpulse* impulse){
    std::cout<< "Postsolve called" << '\n';
  }
  */
};

struct MovementCoordinates {
  MovementCoordinates(b2Vec2 cur,b2Vec2 prev):PreviousPoint(prev),CurrentPoint(cur){}
  b2Vec2 PreviousPoint;
  b2Vec2 CurrentPoint;
};

sf::Vector2f  DisplacementFromPhysicsToPixel(MovementCoordinates movement){
  float x,y;
  x = (movement.CurrentPoint.x - movement.PreviousPoint.x)*SCALE_FACTOR;
  y = -((movement.CurrentPoint.y-HEIGHT) - (movement.PreviousPoint.y-HEIGHT))*SCALE_FACTOR;
  std::cout << "Y disp = "<<y<<'\n';
  return sf::Vector2f(x,y);
}

int main()
{
  //box2d settings
  b2Vec2 gravity(0.0f,-0.05f);
  b2World world(gravity);
  MyContactListener listen;
  world.SetContactListener(&listen);
  float timestep = 1.0f/60.0f;
  int32 velocityIterations = 8;
  int32 positionIterations = 3;

  //sfml settings
  sf::Clock clock;
  float elasped_time =0;

  //ground Configurations
  //box2d configurations

  b2EdgeShape shape;
  shape.Set(
    b2Vec2(0.0f, 0.0f),
    b2Vec2(WIDTH, 0.0f)
    );

  b2FixtureDef sd;
  sd.shape = &shape;
  sd.density = 0.3f;
  sd.friction = 0.0f;
  sd.filter.categoryBits = 0x0002;
  sd.filter.maskBits = 0x0004;

  b2BodyDef bd;
  //bd.position.Set(0.0f,359.0f);
  b2Body* ground = world.CreateBody(&bd);
  ground->CreateFixture(&sd);

  //SFML Declarations
  sf::Vertex line[] = {
    sf::Vertex(sf::Vector2f(0.0f,HEIGHT)),
    sf::Vertex(sf::Vector2f(WIDTH,HEIGHT))
  };
  ///

  //walls
  //box2d declaration
  b2Vec2 wallvertexs[4];
  wallvertexs[0].Set(0.0f,0.0f);
  wallvertexs[1].Set(0.0f,HEIGHT);
  wallvertexs[2].Set(WIDTH,HEIGHT);
  wallvertexs[3].Set(WIDTH,0.0f);
  b2ChainShape wall;
  wall.CreateChain(wallvertexs,4);

  b2FixtureDef wallfixturedef;
  wallfixturedef.shape = &wall;
  wallfixturedef.density =1.0f;
  wallfixturedef.friction=0.0f;

  b2BodyDef wallbodydef;
  //wallbodydef.position.Set(0.0f,0.0f);
  b2Body *wallbody = world.CreateBody(&wallbodydef);
  wallbody->CreateFixture(&wallfixturedef);

  //SFML Declarations
  sf::Vertex walllines[] = {
    sf::Vertex(sf::Vector2f(0.0f,HEIGHT)),
    sf::Vertex(sf::Vector2f(0.0f,0.0f)),
    sf::Vertex(sf::Vector2f(WIDTH,0.0f)),
    sf::Vertex(sf::Vector2f(WIDTH,HEIGHT))
  };

  //ball
  //box2d setup
  b2BodyDef bodyDef;
  bodyDef.type = b2_dynamicBody;
  bodyDef.position.Set(viewSize.x/2,viewSize.y);
  b2Body *body = world.CreateBody(&bodyDef);

  b2CircleShape dynamicCircle;
  dynamicCircle.m_radius = BALL_RADIUS;

  b2FixtureDef fixtureDef;
  fixtureDef.shape = &dynamicCircle;
  fixtureDef.density = 0.3f;
  fixtureDef.friction = 0.0f;
  fixtureDef.filter.categoryBits = 0x0001;
  fixtureDef.filter.maskBits = 0x0006;

  body->CreateFixture(&fixtureDef)->SetRestitution(1.0f);
  //body->SetLinearVelocity(b2Vec2(1.0f,1.0f));

  //SFML declarations
  sf::CircleShape circ(BALL_RADIUS*SCALE_FACTOR);
  circ.setFillColor(sf::Color::White);
  circ.setPosition(WIDTH/2, 0);

  //Paddle
  //Box2d Declarations
  b2BodyDef paddleDef;
  paddleDef.type = b2_dynamicBody;
  paddleDef.position.Set(WIDTH/2,PADDLE_HEIGHT/2);
  //paddleDef.gravityScale = 0.0f;
  b2Body *paddleBody = world.CreateBody(&paddleDef);

  b2PolygonShape paddleShape;
  paddleShape.SetAsBox(PADDLE_WIDTH/2,PADDLE_HEIGHT/2);

  b2FixtureDef paddlefixtureDef;
  paddlefixtureDef.shape = &paddleShape;
  paddlefixtureDef.density = 0.3f;
  paddlefixtureDef.friction = 0.0f;
  paddlefixtureDef.filter.categoryBits = 0x0004;
  paddlefixtureDef.filter.maskBits = 0x0003;

  paddleBody->CreateFixture(&paddlefixtureDef)->SetRestitution(1.0f);

  //SFML Declarations

  sf::RectangleShape paddle(sf::Vector2f(PADDLE_WIDTH*SCALE_FACTOR, PADDLE_HEIGHT*SCALE_FACTOR));
  paddle.setFillColor(sf::Color::White);
  paddle.setOrigin(sf::Vector2f(paddle.getSize().x / 2, paddle.getSize().y / 2));
  paddle.setPosition(WIDTH/2, HEIGHT);

  auto paddle_pos = paddleBody->GetPosition();
  auto new_paddle_pos = paddleBody->GetPosition();
  b2Vec2 ball_pos = body->GetPosition();
  b2Vec2 new_ball_pos = body->GetPosition();

  b2Vec2 grdposition = ground->GetPosition();

  while (window.isOpen()) {
    updateInput();
    window.clear();

    sf::Time delta = clock.restart();
    elasped_time += delta.asSeconds();
    world.Step(timestep,velocityIterations,positionIterations);
    new_ball_pos = body->GetPosition();
    auto ball_displaceemnt = DisplacementFromPhysicsToPixel(MovementCoordinates(new_ball_pos,ball_pos));
    ball_pos = new_ball_pos;
    circ.move(ball_displaceemnt);

    if (playerMoving) {
      b2Vec2 force(getForceVector(delta.asSeconds(),direction).x,0.0);
      paddleBody->ApplyForceToCenter(force,true);
      new_paddle_pos = paddleBody->GetPosition();
      auto paddle_pix_displacement = DisplacementFromPhysicsToPixel(MovementCoordinates(new_paddle_pos,paddle_pos));
      paddle_pos = new_paddle_pos;
      paddle.move(paddle_pix_displacement);
    }else{
      paddleBody->SetLinearVelocity(b2Vec2(0.0f,0.0f));
    }
    window.draw(circ);
    window.draw(line,2,sf::Lines);
    window.draw(walllines,4,sf::Lines);
    window.draw(paddle);
    // Render Game Objects
    window.display();
  }

  std::cout<<"gnd x = "<< grdposition.x <<" gnd y = "<< grdposition.y <<" Paddle x = "<<paddle_pos.x<<" Paddle y = "<<paddle_pos.y<<'\n';
  return 0;
}